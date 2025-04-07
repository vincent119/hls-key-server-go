package configs

// JwtSecret 定義 JWT 密鑰配置
// @Summary JWT secret configuration
// @Description JWT secret configuration
// @Tags JWT
// @ID jwt-secret

type JwtSecret struct {
	Enable      bool   `mapstructure:"enable"`
	SecretKey   string `mapstructure:"secretkey"`
	Expire      int    `mapstructure:"expire"`
	User        string `mapstructure:"user"`
	HeaderKey   string `mapstructure:"header-key"`
	HeaderValue string `mapstructure:"header-value"`
	Iss         string `mapstructure:"iss"`
	Aud         string `mapstructure:"aud"`
}
