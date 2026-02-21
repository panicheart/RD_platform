package config

import (
	"os"
	"strconv"
	"time"

	"rdp-platform/rdp-api/models"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Auth     models.AuthConfig `mapstructure:"auth"`
	Log      LogConfig      `mapstructure:"log"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         string        `mapstructure:"port"`
	Mode         string        `mapstructure:"mode"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	FilePath   string `mapstructure:"file_path"`
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	return &Config{
		Server:   loadServerConfig(),
		Database: loadDatabaseConfig(),
		Auth:     loadAuthConfig(),
		Log:      loadLogConfig(),
	}
}

// loadServerConfig 加载服务器配置
func loadServerConfig() ServerConfig {
	return ServerConfig{
		Host:         getEnv("RDP_SERVER_HOST", "0.0.0.0"),
		Port:         getEnv("RDP_API_PORT", "8080"),
		Mode:         getEnv("RDP_ENV", "development"),
		ReadTimeout:  getDurationEnv("RDP_READ_TIMEOUT", 30*time.Second),
		WriteTimeout: getDurationEnv("RDP_WRITE_TIMEOUT", 30*time.Second),
	}
}

// loadDatabaseConfig 加载数据库配置
func loadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:     getEnv("RDP_DB_HOST", "localhost"),
		Port:     getIntEnv("RDP_DB_PORT", 5432),
		User:     getEnv("RDP_DB_USER", "rdp_user"),
		Password: getEnv("RDP_DB_PASSWORD", "rdp_password"),
		DBName:   getEnv("RDP_DB_NAME", "rdp_db"),
		SSLMode:  getEnv("RDP_DB_SSLMODE", "disable"),
	}
}

// loadAuthConfig 加载认证配置
func loadAuthConfig() models.AuthConfig {
	return models.AuthConfig{
		JWTSecret:       getEnv("RDP_JWT_SECRET", "change-this-secret-in-production"),
		AccessTokenTTL:  getDurationEnv("RDP_ACCESS_TOKEN_TTL", 2*time.Hour),
		RefreshTokenTTL: getDurationEnv("RDP_REFRESH_TOKEN_TTL", 7*24*time.Hour),
		Issuer:          getEnv("RDP_JWT_ISSUER", "rdp-api"),
		Audience:        getEnv("RDP_JWT_AUDIENCE", "rdp-users"),
	}
}

// loadLogConfig 加载日志配置
func loadLogConfig() LogConfig {
	return LogConfig{
		Level:    getEnv("RDP_LOG_LEVEL", "info"),
		Format:   getEnv("RDP_LOG_FORMAT", "json"),
		Output:   getEnv("RDP_LOG_OUTPUT", "stdout"),
		FilePath: getEnv("RDP_LOG_FILE", ""),
	}
}

// DSN 返回数据库连接字符串
func (c *DatabaseConfig) DSN() string {
	return "host=" + c.Host +
		" port=" + strconv.Itoa(c.Port) +
		" user=" + c.User +
		" password=" + c.Password +
		" dbname=" + c.DBName +
		" sslmode=" + c.SSLMode
}

// getEnv 获取环境变量，如果不存在返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getIntEnv 获取整数环境变量
func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getDurationEnv 获取持续时间环境变量
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// IsDevelopment 是否为开发环境
func (c *Config) IsDevelopment() bool {
	return c.Server.Mode == "development"
}

// IsProduction 是否为生产环境
func (c *Config) IsProduction() bool {
	return c.Server.Mode == "production"
}
