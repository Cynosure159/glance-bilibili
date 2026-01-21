// Package main 是程序入口
package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"glance-bilibil/internal/api"
	"glance-bilibil/internal/config"
	"glance-bilibil/internal/service"
)

//go:embed templates/*.html
var templatesFS embed.FS

func main() {
	// 解析命令行参数
	configPath := flag.String("config", "", "配置文件路径")
	port := flag.Int("port", 0, "HTTP 服务端口（覆盖配置文件）")
	flag.Parse()

	// 加载配置
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 命令行端口覆盖配置
	if *port > 0 {
		cfg.Port = *port
	}

	log.Printf("[INFO] 配置加载成功: %d 个 UP 主", len(cfg.Channels))

	// 创建服务
	svc := service.NewVideoService(cfg)

	// 初始化（获取 WBI 密钥等）
	log.Println("[INFO] 正在初始化...")
	if err := svc.Initialize(); err != nil {
		log.Printf("[WARN] 初始化警告: %v (将在首次请求时重试)", err)
	} else {
		log.Println("[INFO] 初始化成功")
	}

	// 创建处理器
	handler, err := api.NewHandler(svc, templatesFS)
	if err != nil {
		log.Fatalf("创建处理器失败: %v", err)
	}

	// 注册路由
	http.HandleFunc("/videos", handler.VideosHandler)
	http.HandleFunc("/health", handler.HealthHandler)
	http.HandleFunc("/", handler.IndexHandler)

	// 启动服务
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("[INFO] 服务启动在 http://localhost%s", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Printf("[ERROR] 服务器错误: %v", err)
		os.Exit(1)
	}
}
