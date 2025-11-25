// Package main provides the entry point for the HLS Key Server
//
//	@title						HLS Key Server API
//	@version					1.0
//	@description				RESTful API for HLS encryption key management with JWT authentication
//	@termsOfService				http://swagger.io/terms/
//
//	@contact.name				API Support
//	@contact.url				https://github.com/vincent119/hls-key-server-go
//	@contact.email				support@example.com
//
//	@license.name				MIT
//	@license.url				https://opensource.org/licenses/MIT
//
//	@host						localhost:8080
//	@BasePath					/api/v1
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				JWT Authorization header using the Bearer scheme. Example: "Bearer {token}"
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	_ "hls-key-server-go/docs"
	"hls-key-server-go/internal/apperrors"
	"hls-key-server-go/internal/configs"
	"hls-key-server-go/internal/handler"
	"hls-key-server-go/internal/repository"
	v1 "hls-key-server-go/internal/routes/api/v1"
	"hls-key-server-go/internal/service"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}
func run() error {
	// Load configuration
	cfg, err := configs.LoadConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// Initialize logger
	logger, err := initLogger(cfg.App.Mode)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer func() {
		// Ignore sync errors on stdout/stderr which are expected
		_ = logger.Sync()
	}()

	// Initialize repository
	keyRepo, err := repository.NewFileKeyRepository("./keys")
	if err != nil {
		return fmt.Errorf("init key repository: %w", err)
	}

	logger.Info("keys loaded",
		zap.Int("count", len(keyRepo.List(context.Background()))),
	)

	// Initialize services
	hlsService := service.NewHLSService(keyRepo, logger)
	authService := service.NewAuthService(&cfg.JwtSecret, logger)

	// Initialize handlers
	hlsHandler := handler.NewHLSHandler(hlsService, logger)
	authHandler := handler.NewAuthHandler(authService, &cfg.JwtSecret, logger)

	// Generate test token for development
	if cfg.App.Mode != "production" {
		token, err := authService.GenerateToken(context.Background(), "test-user")
		if err != nil {
			logger.Error("Failed to generate test token", zap.Error(err))
		} else {
			logger.Info("Test JWT token generated", zap.String("token", token))
			fmt.Println("Test JWT Token:", token)
		}
	}

	// Setup Gin mode
	if strings.EqualFold(cfg.App.Mode, "release") {
		gin.SetMode(gin.ReleaseMode)
	} else if strings.EqualFold(cfg.App.Mode, "debug") {
		gin.SetMode(gin.DebugMode)
	}

	// Create router using new architecture
	router := setupRouter(cfg, hlsHandler, authHandler)

	// Create HTTP server
	serverAddr := ":" + cfg.App.Port
	server := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		logger.Info("server starting", zap.String("addr", serverAddr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server listen error", zap.Error(err))
		}
	}()

	// Setup signal handling
	quit := make(chan os.Signal, 1)
	reload := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	signal.Notify(reload, syscall.SIGHUP)

	// Handle signals
	for {
		select {
		case <-reload:
			// Graceful reload: reload keys without stopping server
			logger.Info("received SIGHUP, reloading keys...")
			if err := hlsService.ReloadKeys(context.Background()); err != nil {
				logger.Error("failed to reload keys", zap.Error(err))
			} else {
				logger.Info("keys reloaded successfully",
					zap.Int("count", len(keyRepo.List(context.Background()))),
				)
			}
		case <-quit:
			// Graceful shutdown
			logger.Info("shutting down server...")

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := server.Shutdown(ctx); err != nil {
				return fmt.Errorf("server forced to shutdown: %w", err)
			}

			logger.Info("server exited")
			return nil
		}
	}
}

func initLogger(mode string) (*zap.Logger, error) {
	var logger *zap.Logger
	var err error

	if strings.EqualFold(mode, "production") {
		logger, err = zap.NewProduction()
	} else {
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger, err = config.Build()
	}

	if err != nil {
		return nil, apperrors.Wrap(err, "create logger")
	}

	return logger, nil
}

// setupRouter creates and configures the Gin router with new handlers
func setupRouter(cfg *configs.Config, hlsHandler *handler.HLSHandler, authHandler *handler.AuthHandler) *gin.Engine {
	// Create Gin instance
	router := gin.New()

	// Setup middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// API v1 routes
	v1Group := router.Group("/api/v1")
	routeGroups := v1.GetRouteGroups(hlsHandler, authHandler)
	for _, routeGroup := range routeGroups {
		routeGroup.RegisterRoutes(v1Group)
	}

	// Health check
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return router
}
