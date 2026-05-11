package main

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
	"microclass-backend/config"
	"microclass-backend/handler"
	"microclass-backend/service"
	"microclass-backend/store"
	"microclass-backend/util"
)

func main() {
	godotenv.Load()

	util.InitLogger("./logs")

	cfg := config.Load()

	s, err := store.New(cfg.DBPath)
	if err != nil {
		log.Fatalf("failed to init store: %v", err)
	}
	defer s.Close()

	mockMode, openAIKey, openAIURL, openAIModel, err := s.GetSettings()
	if err == nil {
		config.InitFromDB(mockMode, openAIKey, openAIURL, openAIModel)
	}

	llm := service.NewLLMService(cfg)
	h := handler.New(s, llm)

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 静态资源（前端 build 产物），非 /api 路径才处理
	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.AbortWithStatus(404)
			return
		}
		// 尝试找到实际文件，找不到则返回 index.html（SPA fallback）
		filePath := "./dist" + c.Request.URL.Path
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			c.File("./dist/index.html")
		} else {
			c.File(filePath)
		}
	})

	api := r.Group("/api")
	{
		api.POST("/upload", h.UploadReference)
		api.POST("/generate", h.Generate)
		api.POST("/generate-stream", h.GenerateStream)
		api.GET("/scripts", h.ListScripts)
		api.GET("/scripts/:id", h.GetScript)
		api.DELETE("/scripts/:id", h.DeleteScript)
		api.GET("/settings", h.GetSettings)
		api.PUT("/settings", h.UpdateSettings)
	}

	log.Printf("Server started on :%s (mock=%v)", cfg.Port, cfg.MockMode)
	log.Printf("日志文件: %s", util.GetLogFilePath())
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
