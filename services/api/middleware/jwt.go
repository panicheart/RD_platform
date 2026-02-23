package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims defines the structure of JWT claims
type JWTClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// JWTMiddleware handles JWT authentication
type JWTMiddleware struct {
	jwtSecret string
	issuer    string
	audience  string
}

// NewJWTMiddleware creates a new JWTMiddleware
func NewJWTMiddleware(jwtSecret, issuer, audience string) *JWTMiddleware {
	return &JWTMiddleware{
		jwtSecret: jwtSecret,
		issuer:    issuer,
		audience:  audience,
	}
}

// GenerateTokenPair generates access and refresh tokens
func (m *JWTMiddleware) GenerateTokenPair(userID, username, role string) (accessToken, refreshToken string, err error) {
	// Access token - 2 hours
	accessClaims := JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    m.issuer,
			Audience:  jwt.ClaimStrings{m.audience},
		},
	}

	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(m.jwtSecret))
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	// Refresh token - 7 days
	refreshClaims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    m.issuer,
		Audience:  jwt.ClaimStrings{m.audience + "/refresh"},
	}

	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(m.jwtSecret))
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

// ValidateToken validates a JWT token and returns claims
func (m *JWTMiddleware) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token claims")
}

// RefreshTokens validates refresh token and generates new token pair
func (m *JWTMiddleware) RefreshTokens(refreshToken string, getUserFunc func(userID string) (username, role string, err error)) (newAccessToken, newRefreshToken string, err error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.jwtSecret), nil
	})

	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["sub"].(string)
		if !ok || userID == "" {
			return "", "", errors.New("invalid refresh token claims")
		}

		username, role, err := getUserFunc(userID)
		if err != nil {
			return "", "", fmt.Errorf("failed to get user: %w", err)
		}

		return m.GenerateTokenPair(userID, username, role)
	}

	return "", "", errors.New("invalid refresh token")
}

// Authenticate returns a Gin middleware that validates JWT tokens
func (m *JWTMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    4010,
				"message": "authorization header required",
				"data":    nil,
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    4011,
				"message": "invalid authorization header format",
				"data":    nil,
			})
			c.Abort()
			return
		}

		claims, err := m.ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    4012,
				"message": "invalid or expired token",
				"data":    nil,
			})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// OptionalAuth returns a Gin middleware that validates JWT tokens if present
func (m *JWTMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.Next()
			return
		}

		claims, err := m.ValidateToken(parts[1])
		if err != nil {
			c.Next()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}
