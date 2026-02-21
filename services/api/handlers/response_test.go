package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return r
}

func TestHealthCheck(t *testing.T) {
	router := setupTestRouter()
	router.GET("/health", HealthCheck)

	t.Run("should return 200 and healthy status", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/health", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response HealthResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "ok", response.Status)
		assert.Equal(t, "1.0.0", response.Version)
	})
}

func TestSuccessResponse(t *testing.T) {
	router := setupTestRouter()
	router.GET("/test", func(c *gin.Context) {
		SuccessResponse(c, map[string]string{"key": "value"})
	})

	t.Run("should return success response", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, 200, response.Code)
		assert.Equal(t, "success", response.Message)
	})
}

func TestErrorResponse(t *testing.T) {
	router := setupTestRouter()
	router.GET("/bad-request", func(c *gin.Context) {
		BadRequestResponse(c, "invalid input")
	})

	t.Run("should return bad request error", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/bad-request", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, 400, response.Code)
		assert.Equal(t, "invalid input", response.Message)
		assert.Nil(t, response.Data)
	})
}

func TestNotFoundResponse(t *testing.T) {
	router := setupTestRouter()
	router.GET("/not-found", func(c *gin.Context) {
		NotFoundResponse(c, "resource not found")
	})

	t.Run("should return not found error", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/not-found", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, 404, response.Code)
	})
}
