package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"rdp/services/api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProjectService is a mock implementation of ProjectService
type MockProjectService struct {
	mock.Mock
}

func (m *MockProjectService) ListProjects(ctx interface{}, page, pageSize int, filters map[string]interface{}) ([]models.Project, int64, error) {
	args := m.Called(ctx, page, pageSize, filters)
	return args.Get(0).([]models.Project), args.Get(1).(int64), args.Error(2)
}

func (m *MockProjectService) GetProjectByID(ctx interface{}, id string) (*models.Project, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Project), args.Error(1)
}

func (m *MockProjectService) CreateProject(ctx interface{}, project *models.Project, userID string) error {
	args := m.Called(ctx, project, userID)
	return args.Error(0)
}

func (m *MockProjectService) UpdateProject(ctx interface{}, id string, updates map[string]interface{}, userID string) (*models.Project, error) {
	args := m.Called(ctx, id, updates, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Project), args.Error(1)
}

func (m *MockProjectService) DeleteProject(ctx interface{}, id, userID string) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func (m *MockProjectService) AddProjectMember(ctx interface{}, projectID, userID, role, addedBy string) (*models.ProjectMember, error) {
	args := m.Called(ctx, projectID, userID, role, addedBy)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ProjectMember), args.Error(1)
}

func (m *MockProjectService) GetProjectMembers(ctx interface{}, projectID string) ([]models.ProjectMember, error) {
	args := m.Called(ctx, projectID)
	return args.Get(0).([]models.ProjectMember), args.Error(1)
}

func (m *MockProjectService) RemoveMember(ctx interface{}, projectID, userID, removedBy string) error {
	args := m.Called(ctx, projectID, userID, removedBy)
	return args.Error(0)
}

func (m *MockProjectService) UpdateMemberRole(ctx interface{}, projectID, userID, newRole, updatedBy string) error {
	args := m.Called(ctx, projectID, userID, newRole, updatedBy)
	return args.Error(0)
}

func (m *MockProjectService) UpdateProjectProgress(ctx interface{}, projectID string, progress int, userID string) (*models.Project, error) {
	args := m.Called(ctx, projectID, progress, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Project), args.Error(1)
}

func (m *MockProjectService) GetUserProjects(ctx interface{}, userID string) ([]models.Project, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Project), args.Error(1)
}

func (m *MockProjectService) GetProjectActivities(ctx interface{}, projectID string) ([]models.Activity, error) {
	args := m.Called(ctx, projectID)
	return args.Get(0).([]models.Activity), args.Error(1)
}

func (m *MockProjectService) CreateActivity(ctx interface{}, activity *models.Activity) error {
	args := m.Called(ctx, activity)
	return args.Error(0)
}

func (m *MockProjectService) GetProjectStats(ctx interface{}, userID string) (map[string]interface{}, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func setupTestRouter() (*gin.Engine, *MockProjectService, *ProjectHandler) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockService := new(MockProjectService)
	handler := NewProjectHandler(mockService)
	return router, mockService, handler
}

func TestProjectHandler_CreateProject(t *testing.T) {
	router, mockService, handler := setupTestRouter()

	router.POST("/api/v1/projects", func(c *gin.Context) {
		c.Set("user_id", "user-123")
		handler.CreateProject(c)
	})

	t.Run("create project successfully", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name":     "Test Project",
			"category": "pd_project",
		}
		jsonBody, _ := json.Marshal(reqBody)

		mockService.On("CreateProject", mock.Anything, mock.AnythingOfType("*models.Project"), "user-123").
			Return(nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/projects", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(0), response["code"])
		assert.Equal(t, "project created successfully", response["message"])
	})

	t.Run("invalid request body", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/projects", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("missing required fields", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"description": "Missing name",
		}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/projects", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestProjectHandler_GetProjects(t *testing.T) {
	router, mockService, handler := setupTestRouter()
	router.GET("/api/v1/projects", handler.GetProjects)

	t.Run("get projects list", func(t *testing.T) {
		projects := []models.Project{
			{ID: uuid.New(), Name: "Project 1", Code: "RDP-PD-20240221-001"},
			{ID: uuid.New(), Name: "Project 2", Code: "RDP-TR-20240221-002"},
		}

		mockService.On("ListProjects", mock.Anything, 1, 20, mock.Anything).
			Return(projects, int64(2), nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/projects?page=1&page_size=20", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(0), response["code"])

		data := response["data"].(map[string]interface{})
		assert.Equal(t, float64(2), data["total"])
	})

	t.Run("get projects with filters", func(t *testing.T) {
		projects := []models.Project{
			{ID: uuid.New(), Name: "Project 1", Code: "RDP-PD-20240221-001", Status: "in_progress"},
		}

		mockService.On("ListProjects", mock.Anything, 1, 20, mock.MatchedBy(func(f map[string]interface{}) bool {
			return f["status"] == "in_progress" && f["category"] == "pd_project"
		})).Return(projects, int64(1), nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/projects?status=in_progress&category=pd_project", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestProjectHandler_GetProject(t *testing.T) {
	router, mockService, handler := setupTestRouter()
	router.GET("/api/v1/projects/:id", handler.GetProject)

	t.Run("get project by id", func(t *testing.T) {
		projectID := uuid.New()
		project := &models.Project{
			ID:       projectID,
			Name:     "Test Project",
			Code:     "RDP-PD-20240221-001",
			Status:   "in_progress",
			Progress: 50,
		}

		mockService.On("GetProjectByID", mock.Anything, projectID.String()).
			Return(project, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/projects/"+projectID.String(), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(0), response["code"])

		data := response["data"].(map[string]interface{})
		assert.Equal(t, "Test Project", data["name"])
	})

	t.Run("project not found", func(t *testing.T) {
		projectID := uuid.New()

		mockService.On("GetProjectByID", mock.Anything, projectID.String()).
			Return(nil, errors.New("project not found")).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/projects/"+projectID.String(), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestProjectHandler_UpdateProject(t *testing.T) {
	router, mockService, handler := setupTestRouter()
	router.PUT("/api/v1/projects/:id", func(c *gin.Context) {
		c.Set("user_id", "user-123")
		handler.UpdateProject(c)
	})

	t.Run("update project successfully", func(t *testing.T) {
		projectID := uuid.New()
		updatedProject := &models.Project{
			ID:   projectID,
			Name: "Updated Project Name",
			Code: "RDP-PD-20240221-001",
		}

		mockService.On("UpdateProject", mock.Anything, projectID.String(), mock.Anything, "user-123").
			Return(updatedProject, nil).Once()

		reqBody := map[string]interface{}{
			"name": "Updated Project Name",
		}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/projects/"+projectID.String(), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(0), response["code"])
	})

	t.Run("insufficient permissions", func(t *testing.T) {
		projectID := uuid.New()

		mockService.On("UpdateProject", mock.Anything, projectID.String(), mock.Anything, "user-123").
			Return(nil, errors.New("insufficient permissions to update project")).Once()

		reqBody := map[string]interface{}{
			"name": "Updated Name",
		}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/projects/"+projectID.String(), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func TestProjectHandler_DeleteProject(t *testing.T) {
	router, mockService, handler := setupTestRouter()
	router.DELETE("/api/v1/projects/:id", func(c *gin.Context) {
		c.Set("user_id", "user-123")
		handler.DeleteProject(c)
	})

	t.Run("delete project successfully", func(t *testing.T) {
		projectID := uuid.New()

		mockService.On("DeleteProject", mock.Anything, projectID.String(), "user-123").
			Return(nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/projects/"+projectID.String(), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(0), response["code"])
	})

	t.Run("project not found", func(t *testing.T) {
		projectID := uuid.New()

		mockService.On("DeleteProject", mock.Anything, projectID.String(), "user-123").
			Return(errors.New("project not found")).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/projects/"+projectID.String(), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestProjectHandler_AddMember(t *testing.T) {
	router, mockService, handler := setupTestRouter()
	router.POST("/api/v1/projects/:id/members", func(c *gin.Context) {
		c.Set("user_id", "user-123")
		handler.AddMember(c)
	})

	t.Run("add member successfully", func(t *testing.T) {
		projectID := uuid.New()
		userID := uuid.New()

		member := &models.ProjectMember{
			ID:        uuid.New(),
			ProjectID: projectID,
			UserID:    userID,
			Role:      "developer",
		}

		mockService.On("AddProjectMember", mock.Anything, projectID.String(), userID.String(), "developer", "user-123").
			Return(member, nil).Once()

		reqBody := map[string]interface{}{
			"user_id": userID.String(),
			"role":    "developer",
		}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/projects/"+projectID.String()+"/members", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(0), response["code"])
	})

	t.Run("invalid request body", func(t *testing.T) {
		projectID := uuid.New()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/projects/"+projectID.String()+"/members", bytes.NewBuffer([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestProjectHandler_GetMembers(t *testing.T) {
	router, mockService, handler := setupTestRouter()
	router.GET("/api/v1/projects/:id/members", handler.GetMembers)

	t.Run("get members successfully", func(t *testing.T) {
		projectID := uuid.New()
		userID1 := uuid.New()
		userID2 := uuid.New()

		members := []models.ProjectMember{
			{ID: uuid.New(), ProjectID: projectID, UserID: userID1, Role: "manager"},
			{ID: uuid.New(), ProjectID: projectID, UserID: userID2, Role: "developer"},
		}

		mockService.On("GetProjectMembers", mock.Anything, projectID.String()).
			Return(members, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/projects/"+projectID.String()+"/members", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(0), response["code"])

		data := response["data"].([]interface{})
		assert.Len(t, data, 2)
	})
}

func TestProjectHandler_UpdateProgress(t *testing.T) {
	router, mockService, handler := setupTestRouter()
	router.PUT("/api/v1/projects/:id/progress", func(c *gin.Context) {
		c.Set("user_id", "user-123")
		handler.UpdateProgress(c)
	})

	t.Run("update progress successfully", func(t *testing.T) {
		projectID := uuid.New()
		updatedProject := &models.Project{
			ID:       projectID,
			Name:     "Test Project",
			Progress: 75,
			Status:   "in_progress",
		}

		mockService.On("UpdateProjectProgress", mock.Anything, projectID.String(), 75, "user-123").
			Return(updatedProject, nil).Once()

		reqBody := map[string]interface{}{
			"progress": 75,
		}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/projects/"+projectID.String()+"/progress", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(0), response["code"])
	})

	t.Run("invalid progress value", func(t *testing.T) {
		projectID := uuid.New()

		reqBody := map[string]interface{}{
			"progress": 150,
		}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/projects/"+projectID.String()+"/progress", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestProjectHandler_GetProjectActivities(t *testing.T) {
	router, mockService, handler := setupTestRouter()
	router.GET("/api/v1/projects/:id/activities", handler.GetProjectActivities)

	t.Run("get activities successfully", func(t *testing.T) {
		projectID := uuid.New()

		activities := []models.Activity{
			{ID: uuid.New(), ProjectID: projectID, Name: "Activity 1", Status: "completed", SortOrder: 1},
			{ID: uuid.New(), ProjectID: projectID, Name: "Activity 2", Status: "in_progress", SortOrder: 2},
		}

		mockService.On("GetProjectActivities", mock.Anything, projectID.String()).
			Return(activities, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/projects/"+projectID.String()+"/activities", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(0), response["code"])

		data := response["data"].([]interface{})
		assert.Len(t, data, 2)
	})
}

func TestProjectHandler_GetUserProjects(t *testing.T) {
	router, mockService, handler := setupTestRouter()
	router.GET("/api/v1/users/me/projects", func(c *gin.Context) {
		c.Set("user_id", "user-123")
		handler.GetUserProjects(c)
	})

	t.Run("get user projects successfully", func(t *testing.T) {
		projects := []models.Project{
			{ID: uuid.New(), Name: "Project 1", Code: "RDP-PD-20240221-001"},
			{ID: uuid.New(), Name: "Project 2", Code: "RDP-TR-20240221-002"},
		}

		mockService.On("GetUserProjects", mock.Anything, "user-123").
			Return(projects, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/users/me/projects", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(0), response["code"])

		data := response["data"].([]interface{})
		assert.Len(t, data, 2)
	})

	t.Run("unauthorized", func(t *testing.T) {
		router2, _, handler2 := setupTestRouter()
		router2.GET("/api/v1/users/me/projects", handler2.GetUserProjects)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/users/me/projects", nil)
		router2.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestProjectHandler_GetProjectStats(t *testing.T) {
	router, mockService, handler := setupTestRouter()
	router.GET("/api/v1/projects/stats", func(c *gin.Context) {
		c.Set("user_id", "user-123")
		handler.GetProjectStats(c)
	})

	t.Run("get stats successfully", func(t *testing.T) {
		stats := map[string]interface{}{
			"total":       int64(100),
			"my_projects": int64(5),
			"by_status": map[string]int64{
				"draft":       10,
				"in_progress": 60,
				"completed":   30,
			},
		}

		mockService.On("GetProjectStats", mock.Anything, "user-123").
			Return(stats, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/projects/stats", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(0), response["code"])

		data := response["data"].(map[string]interface{})
		assert.Equal(t, float64(100), data["total"])
	})
}

func TestProjectHandler_GetProjectGantt(t *testing.T) {
	router, mockService, handler := setupTestRouter()
	router.GET("/api/v1/projects/:id/gantt", handler.GetProjectGantt)

	t.Run("get gantt data successfully", func(t *testing.T) {
		projectID := uuid.New()
		now := time.Now()
		startDate := now
		endDate := now.AddDate(0, 1, 0)

		project := &models.Project{
			ID:       projectID,
			Name:     "Test Project",
			Code:     "RDP-PD-20240221-001",
			Status:   "in_progress",
			Progress: 50,
		}

		activities := []models.Activity{
			{
				ID:        uuid.New(),
				ProjectID: projectID,
				Name:      "Activity 1",
				Status:    "completed",
				Progress:  100,
				StartDate: &startDate,
				EndDate:   &endDate,
				SortOrder: 1,
			},
		}

		mockService.On("GetProjectByID", mock.Anything, projectID.String()).
			Return(project, nil).Once()
		mockService.On("GetProjectActivities", mock.Anything, projectID.String()).
			Return(activities, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/projects/"+projectID.String()+"/gantt", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(0), response["code"])

		data := response["data"].(map[string]interface{})
		assert.NotNil(t, data["project"])
		assert.NotNil(t, data["tasks"])
	})

	t.Run("project not found", func(t *testing.T) {
		projectID := uuid.New()

		mockService.On("GetProjectByID", mock.Anything, projectID.String()).
			Return(nil, errors.New("project not found")).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/projects/"+projectID.String()+"/gantt", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
