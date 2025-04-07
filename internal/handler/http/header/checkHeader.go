package header

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"hls-key-server-go/internal/handler/logging"
)

func Get_header(c *gin.Context, ConfigHeaderString string, ConfigHeaderValue string) bool {
	headerValue := c.GetHeader(ConfigHeaderString)
	if headerValue == "" {
		logging.InitZapLogging().Error("http header fail", zap.String("category", "postmark"))
		return false
	} else if headerValue != ConfigHeaderValue {
		logging.InitZapLogging().Error("http header value fail", zap.String("category", "postmark"))
		return false
	}
	return true
}
