package main

import (
	"log"
	"os"

	"drone-patrol-backend/internal/config"
	"drone-patrol-backend/internal/database"
	"drone-patrol-backend/internal/handlers"
	"drone-patrol-backend/internal/middleware"
	"drone-patrol-backend/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	db, err := database.InitDB(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// 初始化服务
	deviceService := services.NewDeviceService(db)
	mqttService := services.NewMQTTService(db)
	redisService := services.NewRedisService()
	errorCodeService := services.NewErrorCodeService()
	mqttProxy := services.NewMQTTProxyService()
	cameraService := services.NewCameraService(db.DB)

	// 初始化摄像头表
	if err := cameraService.CreateCameraTable(); err != nil {
		log.Printf("Failed to create camera table: %v", err)
	}

	// 插入默认摄像头数据
	if err := cameraService.InsertDefaultCameras(); err != nil {
		log.Printf("Failed to insert default cameras: %v", err)
	}

	// 初始化处理器
	handlers := handlers.NewHandlers(deviceService, mqttService, redisService, errorCodeService, mqttProxy, cameraService)

	// 设置Gin模式
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建路由
	r := gin.Default()

	// 添加中间件
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	// 注册路由
	handlers.RegisterRoutes(r)

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "18080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
