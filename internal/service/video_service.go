// Package service 提供业务逻辑
package service

import (
	"math/rand"
	"sync"
	"time"

	"glance-bilibili/internal/config"
	"glance-bilibili/internal/logger"
	"glance-bilibili/internal/models"
	"glance-bilibili/internal/platform"
	"glance-bilibili/internal/worker"
)

const (
	defaultWorkerCount    = 4
	requestJitterMinDelay = 250 * time.Millisecond
	requestJitterMaxDelay = 1200 * time.Millisecond
)

// cacheEntry 缓存条目
type cacheEntry struct {
	videos    models.VideoList
	updatedAt time.Time
}

// VideoService 视频服务
type VideoService struct {
	client     *platform.BilibiliClient
	config     *config.Config
	cache      map[string]cacheEntry
	mu         sync.RWMutex
	workerPool *worker.Pool
}

// NewVideoService 创建视频服务
func NewVideoService(cfg *config.Config) *VideoService {
	client := platform.NewBilibiliClient()

	// 创建 Worker Pool，降低并发以减少被风控拦截的概率。
	pool := worker.NewPool(defaultWorkerCount)
	pool.Start()

	return &VideoService{
		client:     client,
		config:     cfg,
		cache:      make(map[string]cacheEntry),
		workerPool: pool,
	}
}

// Initialize 初始化服务
func (s *VideoService) Initialize() error {
	return s.client.Initialize()
}

func (s *VideoService) getCachedVideos(mid string, cacheTTLSeconds int) (models.VideoList, bool) {
	// 读取缓存条目
	s.mu.RLock()
	entry, exists := s.cache[mid]
	s.mu.RUnlock()

	if !exists {
		return nil, false
	}

	// 缓存存在但已过期，仍返回旧数据供降级兜底使用
	if time.Since(entry.updatedAt) >= time.Duration(cacheTTLSeconds)*time.Second {
		return entry.videos, false
	}

	return entry.videos, true
}

func (s *VideoService) setCachedVideos(mid string, videos models.VideoList) {
	// 更新缓存
	s.mu.Lock()
	s.cache[mid] = cacheEntry{
		videos:    videos,
		updatedAt: time.Now(),
	}
	s.mu.Unlock()
}

// fetchTask 获取单个频道视频的任务
type fetchTask struct {
	service         *VideoService
	channel         config.ChannelInfo
	limit           int
	cacheTTLSeconds int
	resultChan      chan<- models.VideoList
	wg              *sync.WaitGroup
}

// Execute 实现 worker.Task 接口
func (t *fetchTask) Execute() error {
	defer t.wg.Done()

	// 1. 尝试从缓存获取
	cachedVideos, cacheValid := t.service.getCachedVideos(t.channel.Mid, t.cacheTTLSeconds)
	if cacheValid {
		logger.Debugw("命中有效缓存",
			"up_name", t.channel.Name,
			"up_mid", t.channel.Mid,
			"cached", true,
		)
		t.resultChan <- cachedVideos
		return nil
	}

	// 为非缓存请求增加轻微抖动，避免多个频道同时触发风控。
	time.Sleep(randomRequestDelay())

	// 2. 缓存不存在或已过期，从 API 获取
	videos, err := t.service.client.FetchUserVideos(t.channel.Mid, t.limit, t.channel.Name)
	if err != nil {
		logger.Warnw("获取视频失败",
			"up_name", t.channel.Name,
			"up_mid", t.channel.Mid,
			"error", err,
		)
		// 容错降级：如果 API 失败且有旧缓存，返回旧缓存
		if cachedVideos != nil {
			logger.Infow("API 失败，返回过期缓存数据",
				"up_name", t.channel.Name,
				"cached", true,
			)
			t.resultChan <- cachedVideos
		}
		return err
	}

	// 3. 更新缓存
	t.service.setCachedVideos(t.channel.Mid, videos)

	logger.Infow("获取视频成功",
		"up_name", t.channel.Name,
		"up_mid", t.channel.Mid,
		"video_count", len(videos),
		"cached", false,
	)
	t.resultChan <- videos
	return nil
}

// FetchAllVideos 并发获取所有 UP 主的视频并按时间排序
// cacheTTLSeconds 缓存有效期（秒）
func (s *VideoService) FetchAllVideos(limit int, cacheTTLSeconds int) (models.VideoList, error) {
	if len(s.config.Channels) == 0 {
		return models.VideoList{}, nil
	}

	// 创建结果通道和同步等待组
	videoChan := make(chan models.VideoList, len(s.config.Channels))
	var executeWg sync.WaitGroup

	// 提交任务到 Worker Pool
	for _, channel := range s.config.Channels {
		executeWg.Add(1)

		s.workerPool.Submit(&fetchTask{
			service:         s,
			channel:         channel,
			limit:           limit,
			cacheTTLSeconds: cacheTTLSeconds,
			resultChan:      videoChan,
			wg:              &executeWg,
		})
	}

	// 等待所有任务执行完成后关闭通道
	go func() {
		executeWg.Wait()
		close(videoChan)
	}()

	// 收集结果
	var allVideos models.VideoList
	for videos := range videoChan {
		allVideos = append(allVideos, videos...)
	}

	// 检查是否全部失败且无缓存
	if len(allVideos) == 0 {
		return nil, nil
	}

	// 按时间排序并限制数量
	return allVideos.SortByNewest().Limit(limit), nil
}

// FetchChannelVideos 获取单个 UP 主的视频
func (s *VideoService) FetchChannelVideos(mid string, limit int, cacheTTLSeconds int) (models.VideoList, error) {
	// 1. 尝试从缓存获取
	cachedVideos, cacheValid := s.getCachedVideos(mid, cacheTTLSeconds)
	if cacheValid {
		return cachedVideos.SortByNewest().Limit(limit), nil
	}

	// 为非缓存请求增加轻微抖动，避免与批量抓取同时触发风控。
	time.Sleep(randomRequestDelay())

	// 2. 从 API 获取
	videos, err := s.client.FetchUserVideos(mid, limit, "")
	if err != nil {
		if cachedVideos != nil {
			return cachedVideos.SortByNewest().Limit(limit), nil
		}
		return nil, err
	}

	// 3. 更新缓存
	s.setCachedVideos(mid, videos)

	return videos.SortByNewest().Limit(limit), nil
}

// GetConfig 获取配置
func (s *VideoService) GetConfig() *config.Config {
	return s.config
}

// Shutdown 关闭服务（优雅关闭 Worker Pool）
func (s *VideoService) Shutdown() {
	if s.workerPool != nil {
		s.workerPool.Stop()
	}
}

func randomRequestDelay() time.Duration {
	window := requestJitterMaxDelay - requestJitterMinDelay
	if window <= 0 {
		return requestJitterMinDelay
	}
	return requestJitterMinDelay + time.Duration(rand.Int63n(int64(window)))
}
