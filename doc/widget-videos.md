# Videos Widget 模块文档

## 概述

`widget-videos. go` 是 Glance 应用中用于显示 YouTube 视频内容的小部件模块。该模块实现了从 YouTube 频道或播放列表自动抓取最新上传视频的功能，并支持多种展示样式。

---

## 核心功能

### 1. 视频获取与处理
- **数据源**：通过 YouTube RSS Feed API 获取频道或播放列表中的视频
- **功能特性**：
  - 支持多个 YouTube 频道同时聚合
  - 支持 YouTube 播放列表
  - 可配置是否包含 Shorts 短视频
  - 自动排序（按发布时间最新优先）
  - 支持自定义视频 URL 模板

### 2. 多种展示样式
该模块支持三种不同的视频展示样式：

| 样式 | 模板文件 | 说明 |
|------|--------|------|
| 默认样式 | `videos.html` | 水平卡片轮播展示 |
| `grid-cards` | `videos-grid.html` | 网格卡片布局 |
| `vertical-list` | `videos-vertical-list.html` | 垂直列表展示 |

### 3. 缓存机制
- 默认缓存时长：**1 小时**
- 缓存后无需重复请求，直到过期时自动更新

---

## 数据结构

### videosWidget 结构体
```go
type videosWidget struct {
    widgetBase                  // 继承基础部件属性
    Videos            videoList // 存储视频列表
    VideoUrlTemplate  string    // 自定义URL模板
    Style             string    // 展示风格（grid-cards/vertical-list）
    CollapseAfter     int       // 折叠后显示的项目数
    CollapseAfterRows int       // 网格折叠时的行数（默认4行）
    Channels          []string  // YouTube 频道ID列表
    Playlists         []string  // YouTube 播放列表ID列表
    Limit             int       // 最多显示的视频数（默认25条）
    IncludeShorts     bool      // 是否包含 Shorts 短视频
}
```

### video 结构体
单个视频的数据结构：

```go
type video struct {
    ThumbnailUrl string    // 缩略图URL
    Title        string    // 视频标题
    Url          string    // 视频链接
    Author       string    // 频道名称
    AuthorUrl    string    // 频道URL
    TimePosted   time.Time // 发布时间
}
```

---

## 核心方法

### 1. initialize() 
**功能**：初始化小部件配置

**处理流程**：
- 设置默认标题为 "Videos"
- 设置缓存时长为 1 小时
- 为各配置项设置默认值：
  - `Limit`: 默认 25 个视频
  - `CollapseAfterRows`: 默认 4 行
  - `CollapseAfter`: 默认 7 项
- **关键处理**：将 `Playlists` 参数自动转换为带 `playlist: ` 前缀的频道列表，使用户配置更直观

### 2. update(ctx context.Context)
**功能**：从 YouTube 获取最新视频数据

**执行步骤**：
1. 调用 `fetchYoutubeChannelUploads()` 获取所有频道/播放列表的视频
2. 错误处理和状态检查
3. 根据 `Limit` 限制视频数量
4. 存储处理后的视频列表

### 3. Render() template. HTML
**功能**：根据配置的样式渲染 HTML 内容

**样式匹配**：
```
Style = "grid-cards"      → videos-grid.html
Style = "vertical-list"   → videos-vertical-list.html
默认                       → videos.html
```

### 4. fetchYoutubeChannelUploads()
**功能**：核心数据获取函数，处理 YouTube RSS Feed 请求

**关键特性**：
- **并行请求**：使用 Worker Pool（30个工作线程）并发请求多个频道
- **播放列表处理**：
  ```
  如果是播放列表ID    → 使用 playlist_id 参数
  如果是频道ID且不包含Shorts → 使用 UULF 播放列表ID（仅含完整视频）
  其他频道ID          → 使用 channel_id 参数
  ```
- **URL 模板支持**：支持自定义 URL 模板，使用 `{VIDEO-ID}` 占位符
- **排序**：返回的视频按最新发布时间优先排序
- **错误处理**：
  - 返回 `errPartialContent` 如果部分频道获取失败
  - 返回 `errNoContent` 如果没有获取到任何视频

---

## YouTube RSS Feed 解析

### youtubeFeedResponseXml 结构体
解析 YouTube RSS 源中的数据：

```xml
<feed>
  <author>
    <name>频道名称</name>
    <uri>频道链接</uri>
  </author>
  <entry>  <!-- 单个视频项 -->
    <title>视频标题</title>
    <published>2024-01-16T12:00:00+00:00</published>
    <link href="视频链接" />
    <media:group>
      <media:thumbnail url="缩略图URL" />
    </media:group>
  </entry>
</feed>
```

### 时间解析
支持 YouTube RFC 3339 时间格式：`2006-01-02T15:04:05-07:00`

---

## 配置示例

### 基础配置
```yaml
- type: videos
  title: "最新视频"
  channels:
    - "UC1234567890abcdef"  # YouTube 频道ID
    - "UC9876543210fedcba"
  limit: 20
  style: "grid-cards"
```

### 高级配置
```yaml
- type: videos
  channels:
    - "UC1234567890abcdef"
  playlists:
    - "PLxxxxxxxxxxxx"  # 播放列表ID
  limit: 30
  include-shorts: false
  style: "vertical-list"
  collapse-after: 10
  collapse-after-rows: 3
  video-url-template: "https://custom-domain.com/watch?v={VIDEO-ID}"
```

---

## 处理流程图

```
用户配置 (channels/playlists)
    ↓
initialize() 初始化配置
    ↓
update() 触发数据更新
    ↓
fetchYoutubeChannelUploads()
    ├─ 构建 YouTube RSS 源 URL
    ├─ 并发请求 30 个工作线程
    ├─ 解析 XML 响应
    ├─ 提取视频信息
    └─ 按时间倒序排列
    ↓
Render() 选择模板渲染
    ├─ grid-cards    → 网格卡片
    ├─ vertical-list → 垂直列表
    └─ 默认          → 水平轮播
    ↓
缓存 1 小时后再次更新
```

---

## 技术特点

### 性能优化
- **并发请求**：30 个工作线程并发获取多个频道数据，大幅提升效率
- **缓存策略**：1 小时缓存避免频繁请求 YouTube
- **流式 XML 解析**：直接解析 XML 响应，无需中间转换

### 错误处理
- **部分失败容错**：部分频道获取失败时仍返回已成功获取的视频
- **日志记录**：失败的频道会记录到系统日志便于调试
- **友好的错误返回**：区分完全失败和部分失败情况

### 灵活性
- **多样式支持**：满足不同的界面设计需求
- **自定义 URL**：支持重定向到自建代理或修改 URL 格式
- **Shorts 控制**：可选择是否显示短视频内容

---

## 外部依赖

| 依赖 | 用途 |
|------|------|
| `net/http` | HTTP 请求客户端 |
| `encoding/xml` | RSS Feed XML 解析 |
| `time` | 时间处理和排序 |
| `html/template` | HTML 模板渲染 |
| `net/url` | URL 解析和参数提取 |

---

## 集成方式

该模块在 Glance 框架中的角色：

1. **Widget 接口实现**：实现了 `widget` 接口的所有必要方法
2. **在 widget. go 中注册**：通过 `newWidget("videos")` 创建实例
3. **与框架集成**：
   - 继承 `widgetBase` 获取基础功能
   - 使用框架提供的模板解析器
   - 利用框架的 Worker Pool 机制并发处理

---

## 限制与注意事项

1. **无实时更新**：需要刷新页面才能获取最新视频，默认缓存 1 小时
2. **YouTube API 依赖**：依赖 YouTube RSS Feed，如被限制则功能受影响
3. **频道 ID 格式**：需要使用真实的 YouTube 频道 ID（如 `UCxxx`）
4. **Shorts 独立播放列表**：不包含 Shorts 时会使用特殊的 UULF 播放列表 ID

---

## 总结

`widget-videos.go` 是一个功能完整、性能优化的 YouTube 视频聚合模块，通过：
- 高效的并发数据获取
- 灵活的多样式展示
- 智能的缓存机制
- 人性化的配置设计

为 Glance 用户提供了方便的视频内容集成能力。