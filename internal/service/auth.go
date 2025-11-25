package service

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"

	"hls-key-server-go/internal/apperrors"
	"hls-key-server-go/internal/configs"
)

// AuthService handles authentication logic
type AuthService struct {
	config *configs.JwtSecret
	logger *zap.Logger
}

// NewAuthService creates a new auth service
func NewAuthService(config *configs.JwtSecret, logger *zap.Logger) *AuthService {
	return &AuthService{
		config: config,
		logger: logger,
	}
}

// GenerateToken generates a JWT token for the given username
func (s *AuthService) GenerateToken(_ context.Context, username string) (string, error) {
	claims := jwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(time.Minute * time.Duration(s.config.Expire)).Unix(),
		"iat": time.Now().Unix(),
		"iss": s.config.Iss,
		"aud": s.config.Aud,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.config.SecretKey))
	if err != nil {
		return "", apperrors.Wrap(err, "sign token")
	}

	s.logger.Info("JWT token generated",
		zap.String("username", username),
	)

	return signedToken, nil
}

// ValidateCredentials validates username and custom header
func (s *AuthService) ValidateCredentials(_ context.Context, username, headerValue string) error {
	if username == "" {
		return apperrors.ErrInvalidCredentials
	}

	if username != s.config.User {
		return apperrors.ErrInvalidCredentials
	}

	if headerValue == "" {
		return apperrors.ErrMissingHeader
	}

	if headerValue != s.config.HeaderValue {
		return apperrors.ErrMissingHeader
	}

	return nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *AuthService) ValidateToken(_ context.Context, tokenString string) (jwt.MapClaims, error) {
	if tokenString == "" {
		return nil, apperrors.ErrTokenMissing
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperrors.ErrTokenInvalid
		}
		return []byte(s.config.SecretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, apperrors.Wrap(err, "parse token")
	}

	// Validate required claims
	if !s.validateClaims(claims) {
		return nil, apperrors.ErrTokenInvalid
	}

	return claims, nil
}

func (s *AuthService) validateClaims(claims jwt.MapClaims) bool {
	_, subExists := claims["sub"].(string)
	_, expExists := claims["exp"].(float64)
	_, iatExists := claims["iat"].(float64)
	iss, issExists := claims["iss"].(string)
	aud, audExists := claims["aud"].(string)

	if !subExists || !expExists || !iatExists || !issExists || !audExists {
		return false
	}

	if iss != s.config.Iss || aud != s.config.Aud {
		return false
	}

	return true
}
