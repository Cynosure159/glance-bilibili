# glance-bilibil

为 [glance](https://github.com/glanceapp/glance) 开发的 Bilibili 视频汇总展示插件。

## 功能特性

- ✅ **多 UP 主支持**：通过 `config.json` 配置多个感兴趣的 UP 主。
- ✅ **自动汇总**：并发获取所有配置 UP 主的视频，并按发布时间排序。
- ✅ **风控绕过**：实现 WBI 签名、buvid 获取及 dm 参数模拟，稳定获取数据。
- ✅ **多种样式**：支持默认轮播、网格布局和垂直列表。
- ✅ **配置灵活**：支持配置文件及 URL 参数覆盖。

## 快速开始

### 1. 准备配置
创建 `config.json`:
```json
{
  "port": 8082,
  "channels": [
    { "mid": "946974", "name": "影视飓风" },
    { "mid": "163637592", "name": "老师好我叫何同学" }
  ],
  "limit": 25
}
```

### 2. 运行服务
```bash
go build -o glance-bilibil.exe .
.\glance-bilibil.exe -config config.json
```

### 3. 集成到 Glance
在 `glance.yml` 中添加：
```yaml
- type: extension
  url: http://localhost:8082/videos
  allow-potentially-dangerous-html: true
  cache: 5m
```

## 预览

访问 `http://localhost:8082` 可以查看首页说明及示例。

## 开发

本项目采用分层架构：
- `internal/api`: HTTP 路由与处理
- `internal/service`: 业务逻辑、汇总与排序
- `internal/platform`: Bilibili API 客户端与 WBI 签名
- `internal/config`: 配置管理
- `internal/models`: 核心模型
