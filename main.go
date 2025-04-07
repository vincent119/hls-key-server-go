package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"hls-key-server-go/internal/configs"
	"hls-key-server-go/internal/handler/hls"
	"hls-key-server-go/internal/handler/middleware"
	"hls-key-server-go/internal/routes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var g errgroup.Group

func init() {
	configs.Init()
	mode := configs.Conf.App.Mode
	if strings.ToLower(mode) == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	if strings.ToLower(mode) == "debug" {
		gin.SetMode(gin.DebugMode)
	}
	hls.InitKeys()
	token, _ := middleware.GenerateJWT("test-user")
	fmt.Println(token)
}

func main() {
	port := configs.Conf.App.Port
	serverAddr := ":" + port
	router := routes.DefaultRoute()

	server := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	// 啟動 HTTP 服務
	go func() {
		log.Printf("Server is running on %s", serverAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server listen error: %v", err)
		}
	}()

	// 等待中斷信號 (Ctrl+C) 來優雅關閉服務
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// 設置 5 秒超時
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
