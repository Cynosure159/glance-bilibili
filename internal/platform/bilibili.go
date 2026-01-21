// Package platform 提供 Bilibili API 封装
package platform

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"sync"
	"time"

	"glance-bilibil/internal/models"
)

// bilibiliResponse 用于解析 Bilibili API 响应
type bilibiliResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		List struct {
			Vlist []struct {
				Aid         int64  `json:"aid"`
				Bvid        string `json:"bvid"`
				Title       string `json:"title"`
				Pic         string `json:"pic"`
				Author      string `json:"author"`
				Mid         int64  `json:"mid"`
				Created     int64  `json:"created"`
				Length      string `json:"length"`
				Play        int    `json:"play"`
				Description string `json:"description"`
			} `json:"vlist"`
		} `json:"list"`
		Page struct {
			Count int `json:"count"`
			Pn    int `json:"pn"`
			Ps    int `json:"ps"`
		} `json:"page"`
	} `json:"data"`
}

// buvidResponse 用于解析 buvid API 响应
type buvidResponse struct {
	Code int `json:"code"`
	Data struct {
		B3 string `json:"b_3"`
		B4 string `json:"b_4"`
	} `json:"data"`
}

// BilibiliClient Bilibili API 客户端
type BilibiliClient struct {
	httpClient *http.Client
	wbiKeys    *WbiKeys
	buvid3     string
	buvid4     string
	buvidMu    sync.RWMutex
	webidCache map[string]string
	webidMu    sync.RWMutex
}

// NewBilibiliClient 创建新的 Bilibili 客户端
func NewBilibiliClient() *BilibiliClient {
	return &BilibiliClient{
		httpClient: &http.Client{Timeout: 15 * time.Second},
		wbiKeys:    GetWbiKeys(),
		webidCache: make(map[string]string),
	}
}

// Initialize 初始化客户端（获取必要的密钥）
func (c *BilibiliClient) Initialize() error {
	// 获取 WBI 密钥
	if err := c.wbiKeys.Update(c.httpClient); err != nil {
		return fmt.Errorf("获取 WBI 密钥失败: %w", err)
	}

	// 获取 buvid
	if err := c.ensureBuvid(); err != nil {
		return fmt.Errorf("获取 buvid 失败: %w", err)
	}

	return nil
}

// ensureBuvid 确保已获取 buvid
func (c *BilibiliClient) ensureBuvid() error {
	c.buvidMu.RLock()
	hasValue := c.buvid3 != "" && c.buvid4 != ""
	c.buvidMu.RUnlock()

	if hasValue {
		return nil
	}

	c.buvidMu.Lock()
	defer c.buvidMu.Unlock()

	if c.buvid3 != "" && c.buvid4 != "" {
		return nil
	}

	req, err := http.NewRequest("GET", "https://api.bilibili.com/x/frontend/finger/spi", nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Referer", "https://www.bilibili.com/")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var buvidResp buvidResponse
	if err := json.Unmarshal(body, &buvidResp); err != nil {
		return err
	}

	if buvidResp.Code != 0 {
		return fmt.Errorf("获取 buvid 失败: code=%d", buvidResp.Code)
	}

	c.buvid3 = buvidResp.Data.B3
	c.buvid4 = buvidResp.Data.B4
	log.Printf("[INFO] 获取 buvid 成功")

	return nil
}

// getWebid 获取 w_webid 参数
func (c *BilibiliClient) getWebid(mid string) string {
	c.webidMu.RLock()
	if webid, ok := c.webidCache[mid]; ok {
		c.webidMu.RUnlock()
		return webid
	}
	c.webidMu.RUnlock()

	c.webidMu.Lock()
	defer c.webidMu.Unlock()

	if webid, ok := c.webidCache[mid]; ok {
		return webid
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://space.bilibili.com/%s/dynamic", mid), nil)
	if err != nil {
		return ""
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	re := regexp.MustCompile(`"access_id"\s*:\s*"([^"]+)"`)
	matches := re.FindSubmatch(body)
	if len(matches) > 1 {
		webid := string(matches[1])
		c.webidCache[mid] = webid
		return webid
	}

	return ""
}

// getDmParams 生成 dm 相关参数
func getDmParams() url.Values {
	chars := "ABCDEFGHIJK"
	randomChars := func() string {
		result := make([]byte, 2)
		for i := range result {
			result[i] = chars[rand.Intn(len(chars))]
		}
		return string(result)
	}

	params := url.Values{}
	params.Set("dm_img_list", "[]")
	params.Set("dm_img_str", randomChars())
	params.Set("dm_cover_img_str", randomChars())
	params.Set("dm_img_inter", `{"ds":[],"wh":[0,0,0],"of":[0,0,0]}`)
	return params
}

// FetchUserVideos 获取指定用户的视频列表
// authorOverride 如果非空，则用它覆盖 API 返回的作者名称（解决联合投稿问题）
func (c *BilibiliClient) FetchUserVideos(mid string, limit int, authorOverride string) (models.VideoList, error) {
	if err := c.ensureBuvid(); err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Set("mid", mid)
	params.Set("order", "pubdate")
	params.Set("pn", "1")
	params.Set("ps", strconv.Itoa(min(limit, 50)))
	params.Set("jsonp", "jsonp")

	// 添加 dm 参数
	for k, v := range getDmParams() {
		params[k] = v
	}

	// 获取 webid
	if webid := c.getWebid(mid); webid != "" {
		params.Set("w_webid", webid)
	}

	// WBI 签名
	signedParams, err := c.wbiKeys.Sign(params, c.httpClient)
	if err != nil {
		return nil, fmt.Errorf("WBI 签名失败: %w", err)
	}

	apiURL := "https://api.bilibili.com/x/space/wbi/arc/search?" + signedParams.Encode()

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Referer", "https://space.bilibili.com/"+mid)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Origin", "https://space.bilibili.com")

	c.buvidMu.RLock()
	cookie := fmt.Sprintf("buvid3=%s; buvid4=%s", c.buvid3, c.buvid4)
	c.buvidMu.RUnlock()
	req.Header.Set("Cookie", cookie)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp bilibiliResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, err
	}

	if apiResp.Code != 0 {
		return nil, fmt.Errorf("API 错误: code=%d, message=%s", apiResp.Code, apiResp.Message)
	}

	videos := make(models.VideoList, 0, len(apiResp.Data.List.Vlist))
	for _, v := range apiResp.Data.List.Vlist {
		// 使用配置中的名称覆盖（如果提供）
		author := v.Author
		if authorOverride != "" {
			author = authorOverride
		}

		video := models.Video{
			Title:        v.Title,
			ThumbnailUrl: v.Pic,
			Url:          fmt.Sprintf("https://www.bilibili.com/video/%s", v.Bvid),
			Author:       author,
			AuthorUrl:    fmt.Sprintf("https://space.bilibili.com/%s", mid),
			TimePosted:   time.Unix(v.Created, 0),
			Duration:     v.Length,
			PlayCount:    v.Play,
			Bvid:         v.Bvid,
		}
		videos = append(videos, video)
		if len(videos) >= limit {
			break
		}
	}

	return videos, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
