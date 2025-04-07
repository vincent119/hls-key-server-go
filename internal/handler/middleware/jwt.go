package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"hls-key-server-go/internal/configs"
	"hls-key-server-go/internal/handler/logging"
	"log"
	"net/http"
	"strings"
	"time"
)

func getJWTSecret() []byte {
	return []byte(configs.Conf.JwtSecret.SecretKey)
}

func GenerateJWT(username string) (string, error) {
	claims := jwt.MapClaims{
		//"sub": configs.Conf.JwtSecret.User,
		"sub": username,
		"exp": time.Now().Add(time.Minute * time.Duration(configs.Conf.JwtSecret.Expire)).Unix(),
		"iat": time.Now().Unix(),
		"iss": configs.Conf.JwtSecret.Iss,
		"aud": configs.Conf.JwtSecret.Aud,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(getJWTSecret())
	if err != nil {
		return "", err
	}

	logging.InitZapLogging().Info("JWT token generated",
		//zap.String("ip", userIP),
		zap.String("username", username),
		zap.String("token", signedToken),
	)

	return signedToken, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// get token from header, query or cookie
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			tokenString = c.Query("token")
			if tokenString == "" {
				tokenString, _ = c.Cookie("token")
			}
		}

		// if token is empty return 401
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
			c.Abort()

			logging.InitZapLogging().Error("Token is required",
				zap.String("ip", c.ClientIP()),
				//zap.String("ip", userIP),
			)
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return getJWTSecret(), nil
		})

		// check Claims

		if err != nil || !token.Valid {
			log.Println("JWT Invalid or expired token:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		sub, subExists := claims["sub"].(string)
		expFloat, expExists := claims["exp"].(float64)
		iatFloat, iatExists := claims["iat"].(float64)
		iss, issExists := claims["iss"].(string)
		aud, audExists := claims["aud"].(string)

		exp := int64(expFloat)
		iat := int64(iatFloat)

		if !subExists || !expExists || !iatExists || !issExists || !audExists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: missing claims"})
			c.Abort()
			return
		}

		// check  Token expiration
		if time.Now().Unix() > int64(exp) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
			c.Abort()
			return
		}

		if int64(iat) > time.Now().Unix() {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token issued in the future"})
			c.Abort()
			return
		}

		// time expired access denied
		if iss != configs.Conf.JwtSecret.Iss {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token issuer"})
			c.Abort()
			return
		}

		if aud != configs.Conf.JwtSecret.Aud {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token audience"})
			c.Abort()
			return
		}

		//  Record user information for further processing
		c.Set("user", sub)
		c.Next()
	}
}
