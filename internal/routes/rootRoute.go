package routes

import (
	_ "hls-key-server-go/docs"
	"hls-key-server-go/internal/configs"
	//"hls-key-server-go/internal/handler/api"
	"hls-key-server-go/internal/handler/middleware"
	v1 "hls-key-server-go/internal/routes/api/v1"
	"io"
	"log"
	"os"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// 需要跳過日誌的 API 路徑
var skipPathList = []string{
	"/healthz",
	"/api/v1/metrics",
}

// RouteGroup 定義 API 組的接口，用於註冊路由
type RouteGroup interface {
	RegisterRoutes(router *gin.RouterGroup)
}

// DefaultRoute 設置 Gin 路由
// @title HLS Key Server API
// @version 1.0
// @description This is the API documentation for HLS Key Server.
// @host localhost:9090
// @BasePath /
func DefaultRoute() *gin.Engine {
	// 設定 Gin 輸出
	gin.ForceConsoleColor()
	gin.DefaultWriter = io.MultiWriter(os.Stdout)

	// 創建 Gin 實例
	routes := gin.New()

	// 設定中介件
	routes.Use(
		otelgin.Middleware("hls-key-server-go"),                         // OpenTelemetry 追蹤
		middleware.CORS(),                                               // 跨域處理
		gin.LoggerWithConfig(gin.LoggerConfig{SkipPaths: skipPathList}), // 設定日誌跳過路徑
		gzip.Gzip(
			gzip.DefaultCompression,
			gzip.WithExcludedPaths([]string{
				"/healthz",
				"/api/v1/metrics",
				"/swagger/*any",
				"/swagger/doc.json",
			})),
		gin.Recovery(),
	)

	// 設置信任代理 (允許從配置讀取，而不是硬編碼 IP)
	trustedProxies := []string{"172.16.99.200"}
	if configs.Conf.App.Mode == "production" {
		trustedProxies = []string{"192.168.1.1", "10.0.0.1"} // 這裡可以從配置讀取
	}
	if err := routes.SetTrustedProxies(trustedProxies); err != nil {
		log.Fatalf("無法設置信任代理: %v", err)
	}

	// 健康檢查路由
	routes.GET("/healthz", HealthCheck)

	// Prometheus 監控指標
	p := ginprometheus.NewPrometheus("gin")

	// 讀取 Prometheus 用戶密碼
	metricUser := configs.Conf.Metric.User
	metricPassword := configs.Conf.Metric.Password
	if metricUser == "" || metricPassword == "" {
		log.Println("Access Deny")
		p.Use(routes)
	} else {
		p.MetricsPath = "/api/v1/metrics"
		p.UseWithAuth(routes, gin.Accounts{metricUser: metricPassword})
	}

	routes.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routeGroups := v1.GetRouteGroups()
	// routeGroups := []RouteGroup{}

	apiGroup := routes.Group("/api/v1")
	for _, group := range routeGroups {
		group.RegisterRoutes(apiGroup)
	}

	return routes
}
