package models

import "time"

// AuthConfig 认证配置
type AuthConfig struct {
	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	Issuer          string
	Audience        string
}
