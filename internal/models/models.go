// Package models 定义数据结构
package models

import (
	"sort"
	"time"
)

// Video 表示单个视频的信息
type Video struct {
	Title        string    `json:"title"`         // 视频标题
	ThumbnailUrl string    `json:"thumbnail_url"` // 缩略图 URL
	Url          string    `json:"url"`           // 视频链接
	Author       string    `json:"author"`        // UP 主名称
	AuthorUrl    string    `json:"author_url"`    // UP 主主页链接
	TimePosted   time.Time `json:"time_posted"`   // 发布时间
	Duration     string    `json:"duration"`      // 视频时长
	PlayCount    int       `json:"play_count"`    // 播放次数
	Bvid         string    `json:"bvid"`          // BV 号
}

// VideoList 视频列表类型
type VideoList []Video

// SortByNewest 按发布时间倒序排序（最新在前）
func (v VideoList) SortByNewest() VideoList {
	sort.Slice(v, func(i, j int) bool {
		return v[i].TimePosted.After(v[j].TimePosted)
	})
	return v
}

// Limit 限制返回数量
func (v VideoList) Limit(n int) VideoList {
	if n <= 0 || n >= len(v) {
		return v
	}
	return v[:n]
}
