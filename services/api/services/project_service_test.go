package services

import (
	"context"
	"testing"
	"time"

	"rdp/services/api/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open gorm: %v", err)
	}

	return db, mock
}

func TestProjectService_GenerateProjectCode(t *testing.T) {
	db, mock := setupMockDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	service := NewProjectService(db)
	ctx := context.Background()

	t.Run("generate first code of the day", func(t *testing.T) {
		today := time.Now().Format("20060102")
		prefix := "RDP-PD-" + today + "-"

		mock.ExpectQuery(`SELECT code FROM "projects"`).
			WithArgs(prefix + "%").
			WillReturnRows(sqlmock.NewRows([]string{"code"}))

		mock.ExpectQuery(`SELECT count\\(\\*\\) FROM "projects"`).
			WithArgs(prefix+"%03d", 1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		code, err := service.GenerateProjectCode(ctx, "pd_project")
		assert.NoError(t, err)
		assert.Contains(t, code, "RDP-PD-")
		assert.Contains(t, code, today)
	})

	t.Run("generate sequential code", func(t *testing.T) {
		today := time.Now().Format("20060102")
		prefix := "RDP-TR-" + today + "-"
		existingCode := prefix + "005"

		mock.ExpectQuery(`SELECT code FROM "projects"`).
			WithArgs(prefix + "%").
			WillReturnRows(sqlmock.NewRows([]string{"code"}).AddRow(existingCode))

		mock.ExpectQuery(`SELECT count\\(\\*\\) FROM "projects"`).
			WithArgs(prefix+"%03d", 6).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		code, err := service.GenerateProjectCode(ctx, "tech_research")
		assert.NoError(t, err)
		assert.Equal(t, prefix+"006", code)
	})
}

func TestProjectService_GetProjectByID(t *testing.T) {
	db, mock := setupMockDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	service := NewProjectService(db)
	ctx := context.Background()

	t.Run("get existing project", func(t *testing.T) {
		projectID := uuid.New()
		now := time.Now()

		rows := sqlmock.NewRows([]string{"id", "code", "name", "category", "status", "progress", "created_at", "updated_at"}).
			AddRow(projectID, "RDP-PD-20240221-001", "Test Project", "pd_project", "draft", 0, now, now)

		mock.ExpectQuery(`SELECT \* FROM "projects"`).
			WithArgs(projectID).
			WillReturnRows(rows)

		mock.ExpectQuery(`SELECT \* FROM "project_members"`).
			WithArgs(projectID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "project_id", "user_id", "role"}))

		project, err := service.GetProjectByID(ctx, projectID.String())
		assert.NoError(t, err)
		assert.NotNil(t, project)
		assert.Equal(t, "Test Project", project.Name)
	})

	t.Run("get non-existent project", func(t *testing.T) {
		projectID := uuid.New()

		mock.ExpectQuery(`SELECT \* FROM "projects"`).
			WithArgs(projectID).
			WillReturnError(gorm.ErrRecordNotFound)

		project, err := service.GetProjectByID(ctx, projectID.String())
		assert.Error(t, err)
		assert.Nil(t, project)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("invalid project ID", func(t *testing.T) {
		project, err := service.GetProjectByID(ctx, "invalid-id")
		assert.Error(t, err)
		assert.Nil(t, project)
		assert.Contains(t, err.Error(), "invalid")
	})
}

func TestProjectService_ListProjects(t *testing.T) {
	db, mock := setupMockDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	service := NewProjectService(db)
	ctx := context.Background()

	t.Run("list projects with pagination", func(t *testing.T) {
		now := time.Now()
		projectID1 := uuid.New()
		projectID2 := uuid.New()

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)

		projectRows := sqlmock.NewRows([]string{"id", "code", "name", "category", "status", "progress", "created_at", "updated_at"}).
			AddRow(projectID1, "RDP-PD-20240221-001", "Project 1", "pd_project", "draft", 0, now, now).
			AddRow(projectID2, "RDP-TR-20240221-002", "Project 2", "tech_research", "in_progress", 50, now, now)

		mock.ExpectQuery(`SELECT count\(\*\) FROM "projects"`).WillReturnRows(countRows)
		mock.ExpectQuery(`SELECT \* FROM "projects"`).WillReturnRows(projectRows)

		filters := make(map[string]interface{})
		projects, total, err := service.ListProjects(ctx, 1, 20, filters)

		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, projects, 2)
	})

	t.Run("list projects with status filter", func(t *testing.T) {
		now := time.Now()
		projectID := uuid.New()

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)

		projectRows := sqlmock.NewRows([]string{"id", "code", "name", "category", "status", "progress", "created_at", "updated_at"}).
			AddRow(projectID, "RDP-PD-20240221-001", "Project 1", "pd_project", "in_progress", 50, now, now)

		mock.ExpectQuery(`SELECT count\(\*\) FROM "projects"`).
			WithArgs("in_progress").
			WillReturnRows(countRows)
		mock.ExpectQuery(`SELECT \* FROM "projects"`).
			WithArgs("in_progress").
			WillReturnRows(projectRows)

		filters := map[string]interface{}{
			"status": "in_progress",
		}
		projects, total, err := service.ListProjects(ctx, 1, 20, filters)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, projects, 1)
		assert.Equal(t, "in_progress", projects[0].Status)
	})

	t.Run("list projects with search", func(t *testing.T) {
		now := time.Now()
		projectID := uuid.New()

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)

		projectRows := sqlmock.NewRows([]string{"id", "code", "name", "category", "status", "progress", "created_at", "updated_at"}).
			AddRow(projectID, "RDP-PD-20240221-001", "Test Search Project", "pd_project", "draft", 0, now, now)

		mock.ExpectQuery(`SELECT count\(\*\) FROM "projects"`).WillReturnRows(countRows)
		mock.ExpectQuery(`SELECT \* FROM "projects"`).WillReturnRows(projectRows)

		filters := map[string]interface{}{
			"search": "Search",
		}
		projects, total, err := service.ListProjects(ctx, 1, 20, filters)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
	})
}

func TestProjectService_CreateProject(t *testing.T) {
	db, mock := setupMockDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	service := NewProjectService(db)
	ctx := context.Background()
	userID := uuid.New().String()

	t.Run("create project with auto-generated code", func(t *testing.T) {
		today := time.Now().Format("20060102")
		prefix := "RDP-PD-" + today + "-"

		project := &models.Project{
			Name:     "New Project",
			Category: "pd_project",
		}

		// Mock code generation query
		mock.ExpectQuery(`SELECT code FROM "projects"`).WillReturnRows(sqlmock.NewRows([]string{"code"}))
		mock.ExpectQuery(`SELECT count\(\*\) FROM "projects"`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		// Mock transaction
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "projects"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectQuery(`INSERT INTO "project_members"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
		mock.ExpectCommit()

		err := service.CreateProject(ctx, project, userID)
		assert.NoError(t, err)
		assert.NotEmpty(t, project.Code)
		assert.Contains(t, project.Code, prefix)
		assert.Equal(t, "draft", project.Status)
	})

	t.Run("create project with duplicate code", func(t *testing.T) {
		project := &models.Project{
			Name:     "New Project",
			Code:     "RDP-PD-20240221-001",
			Category: "pd_project",
		}

		mock.ExpectQuery(`SELECT count\(\*\) FROM "projects"`).
			WithArgs("RDP-PD-20240221-001").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		err := service.CreateProject(ctx, project, userID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestProjectService_UpdateProject(t *testing.T) {
	db, mock := setupMockDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	service := NewProjectService(db)
	ctx := context.Background()
	userID := uuid.New().String()
	projectID := uuid.New()

	t.Run("update project successfully", func(t *testing.T) {
		now := time.Now()
		updates := map[string]interface{}{
			"name": "Updated Name",
		}

		// Check permission - user is admin
		userRows := sqlmock.NewRows([]string{"id", "role"}).
			AddRow(uuid.MustParse(userID), "admin")

		mock.ExpectQuery(`SELECT \* FROM "users"`).
			WithArgs(sqlmock.AnyArg()).
			WillReturnRows(userRows)

		// Update project
		mock.ExpectExec(`UPDATE "projects" SET`).
			WithArgs("Updated Name", sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Get updated project
		projectRows := sqlmock.NewRows([]string{"id", "code", "name", "category", "status", "created_at", "updated_at"}).
			AddRow(projectID, "RDP-PD-20240221-001", "Updated Name", "pd_project", "draft", now, now)

		mock.ExpectQuery(`SELECT \* FROM "projects"`).
			WithArgs(projectID).
			WillReturnRows(projectRows)

		mock.ExpectQuery(`SELECT \* FROM "project_members"`).
			WithArgs(projectID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "project_id", "user_id", "role"}))

		project, err := service.UpdateProject(ctx, projectID.String(), updates, userID)
		assert.NoError(t, err)
		assert.NotNil(t, project)
		assert.Equal(t, "Updated Name", project.Name)
	})

	t.Run("update non-existent project", func(t *testing.T) {
		updates := map[string]interface{}{
			"name": "Updated Name",
		}

		userRows := sqlmock.NewRows([]string{"id", "role"}).
			AddRow(uuid.MustParse(userID), "admin")

		mock.ExpectQuery(`SELECT \* FROM "users"`).
			WithArgs(sqlmock.AnyArg()).
			WillReturnRows(userRows)

		mock.ExpectExec(`UPDATE "projects" SET`).
			WillReturnResult(sqlmock.NewResult(0, 0))

		project, err := service.UpdateProject(ctx, projectID.String(), updates, userID)
		assert.Error(t, err)
		assert.Nil(t, project)
	})
}

func TestProjectService_DeleteProject(t *testing.T) {
	db, mock := setupMockDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	service := NewProjectService(db)
	ctx := context.Background()
	userID := uuid.New().String()
	projectID := uuid.New()

	t.Run("delete project successfully", func(t *testing.T) {
		// Check permission
		userRows := sqlmock.NewRows([]string{"id", "role"}).
			AddRow(uuid.MustParse(userID), "admin")

		mock.ExpectQuery(`SELECT \* FROM "users"`).
			WithArgs(sqlmock.AnyArg()).
			WillReturnRows(userRows)

		mock.ExpectExec(`UPDATE "projects" SET`).
			WithArgs("deleted", sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := service.DeleteProject(ctx, projectID.String(), userID)
		assert.NoError(t, err)
	})

	t.Run("delete non-existent project", func(t *testing.T) {
		// Check permission
		userRows := sqlmock.NewRows([]string{"id", "role"}).
			AddRow(uuid.MustParse(userID), "admin")

		mock.ExpectQuery(`SELECT \* FROM "users"`).
			WithArgs(sqlmock.AnyArg()).
			WillReturnRows(userRows)

		mock.ExpectExec(`UPDATE "projects" SET`).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := service.DeleteProject(ctx, projectID.String(), userID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestProjectService_AddProjectMember(t *testing.T) {
	db, mock := setupMockDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	service := NewProjectService(db)
	ctx := context.Background()
	projectID := uuid.New().String()
	userID := uuid.New().String()
	addedBy := uuid.New().String()

	t.Run("add member successfully", func(t *testing.T) {
		// Check permission
		userRows := sqlmock.NewRows([]string{"id", "role"}).
			AddRow(uuid.MustParse(addedBy), "admin")

		mock.ExpectQuery(`SELECT \* FROM "users"`).
			WithArgs(sqlmock.AnyArg()).
			WillReturnRows(userRows)

		// Check existing member
		mock.ExpectQuery(`SELECT \* FROM "project_members"`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		// Create member
		mock.ExpectQuery(`INSERT INTO "project_members"`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))

		// Load user info
		memberUserRows := sqlmock.NewRows([]string{"id", "username", "display_name"}).
			AddRow(uuid.MustParse(userID), "testuser", "Test User")

		mock.ExpectQuery(`SELECT \* FROM "users"`).
			WithArgs(sqlmock.AnyArg()).
			WillReturnRows(memberUserRows)

		member, err := service.AddProjectMember(ctx, projectID, userID, "developer", addedBy)
		assert.NoError(t, err)
		assert.NotNil(t, member)
		assert.Equal(t, "developer", member.Role)
	})

	t.Run("add duplicate member", func(t *testing.T) {
		// Check permission
		userRows := sqlmock.NewRows([]string{"id", "role"}).
			AddRow(uuid.MustParse(addedBy), "admin")

		mock.ExpectQuery(`SELECT \* FROM "users"`).
			WithArgs(sqlmock.AnyArg()).
			WillReturnRows(userRows)

		// Check existing member - found
		existingRows := sqlmock.NewRows([]string{"id", "project_id", "user_id", "role"}).
			AddRow(uuid.New(), projectID, userID, "member")

		mock.ExpectQuery(`SELECT \* FROM "project_members"`).
			WillReturnRows(existingRows)

		member, err := service.AddProjectMember(ctx, projectID, userID, "developer", addedBy)
		assert.Error(t, err)
		assert.Nil(t, member)
		assert.Contains(t, err.Error(), "already a member")
	})
}

func TestProjectService_GetProjectMembers(t *testing.T) {
	db, mock := setupMockDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	service := NewProjectService(db)
	ctx := context.Background()
	projectID := uuid.New().String()

	t.Run("get members successfully", func(t *testing.T) {
		memberID1 := uuid.New()
		memberID2 := uuid.New()
		userID1 := uuid.New()
		userID2 := uuid.New()

		memberRows := sqlmock.NewRows([]string{"id", "project_id", "user_id", "role", "joined_at"}).
			AddRow(memberID1, projectID, userID1, "manager", time.Now()).
			AddRow(memberID2, projectID, userID2, "developer", time.Now())

		mock.ExpectQuery(`SELECT \* FROM "project_members"`).
			WillReturnRows(memberRows)

		members, err := service.GetProjectMembers(ctx, projectID)
		assert.NoError(t, err)
		assert.Len(t, members, 2)
	})
}

func TestProjectService_UpdateProjectProgress(t *testing.T) {
	db, mock := setupMockDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	service := NewProjectService(db)
	ctx := context.Background()
	userID := uuid.New().String()
	projectID := uuid.New()

	t.Run("update progress successfully", func(t *testing.T) {
		now := time.Now()

		// Check permission
		userRows := sqlmock.NewRows([]string{"id", "role"}).
			AddRow(uuid.MustParse(userID), "admin")

		mock.ExpectQuery(`SELECT \* FROM "users"`).
			WithArgs(sqlmock.AnyArg()).
			WillReturnRows(userRows)

		// Get activities for calculation
		mock.ExpectQuery(`SELECT \* FROM "activities"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "project_id", "status", "progress"}))

		// Update project
		mock.ExpectExec(`UPDATE "projects" SET`).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Get updated project
		projectRows := sqlmock.NewRows([]string{"id", "code", "name", "category", "status", "progress", "created_at", "updated_at"}).
			AddRow(projectID, "RDP-PD-20240221-001", "Test Project", "pd_project", "in_progress", 50, now, now)

		mock.ExpectQuery(`SELECT \* FROM "projects"`).
			WithArgs(projectID).
			WillReturnRows(projectRows)

		mock.ExpectQuery(`SELECT \* FROM "project_members"`).
			WithArgs(projectID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "project_id", "user_id", "role"}))

		project, err := service.UpdateProjectProgress(ctx, projectID.String(), 50, userID)
		assert.NoError(t, err)
		assert.NotNil(t, project)
	})

	t.Run("invalid progress value", func(t *testing.T) {
		project, err := service.UpdateProjectProgress(ctx, projectID.String(), 150, userID)
		assert.Error(t, err)
		assert.Nil(t, project)
		assert.Contains(t, err.Error(), "between 0 and 100")
	})
}

func TestProjectService_GetProjectStats(t *testing.T) {
	db, mock := setupMockDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	service := NewProjectService(db)
	ctx := context.Background()
	userID := uuid.New().String()

	t.Run("get stats successfully", func(t *testing.T) {
		// Status counts
		statusRows := sqlmock.NewRows([]string{"status", "count"}).
			AddRow("draft", 5).
			AddRow("in_progress", 10).
			AddRow("completed", 3)

		mock.ExpectQuery(`SELECT status, count`).WillReturnRows(statusRows)

		// Category counts
		categoryRows := sqlmock.NewRows([]string{"category", "count"}).
			AddRow("pd_project", 8).
			AddRow("tech_research", 5).
			AddRow("pre_research", 5)

		mock.ExpectQuery(`SELECT category, count`).WillReturnRows(categoryRows)

		// User's projects
		mock.ExpectQuery(`SELECT count\(\*\) FROM "projects"`).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))

		// Total count
		mock.ExpectQuery(`SELECT count\(\*\) FROM "projects"`).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(18))

		stats, err := service.GetProjectStats(ctx, userID)
		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, int64(18), stats["total"])
	})
}

func TestCalculateProgressFromActivities(t *testing.T) {
	t.Run("no activities", func(t *testing.T) {
		activities := []models.Activity{}
		progress := calculateProgressFromActivities(activities)
		assert.Equal(t, 0, progress)
	})

	t.Run("all pending", func(t *testing.T) {
		activities := []models.Activity{
			{Status: "pending", Progress: 0},
			{Status: "pending", Progress: 0},
		}
		progress := calculateProgressFromActivities(activities)
		assert.Equal(t, 0, progress)
	})

	t.Run("all completed", func(t *testing.T) {
		activities := []models.Activity{
			{Status: "completed", Progress: 100},
			{Status: "completed", Progress: 100},
		}
		progress := calculateProgressFromActivities(activities)
		assert.Equal(t, 100, progress)
	})

	t.Run("mixed status", func(t *testing.T) {
		activities := []models.Activity{
			{Status: "completed", Progress: 100},
			{Status: "in_progress", Progress: 50},
			{Status: "pending", Progress: 0},
		}
		progress := calculateProgressFromActivities(activities)
		// (100 + 50 + 0) / 3 = 50
		assert.Equal(t, 50, progress)
	})

	t.Run("partial progress", func(t *testing.T) {
		activities := []models.Activity{
			{Status: "completed", Progress: 100},
			{Status: "in_progress", Progress: 25},
		}
		progress := calculateProgressFromActivities(activities)
		// (100 + 25) / 2 = 62.5 -> 62 (integer division)
		assert.Equal(t, 62, progress)
	})
}

func TestGetCategoryCode(t *testing.T) {
	tests := []struct {
		category string
		expected string
	}{
		{"pd_project", "PD"},
		{"pre_research", "PR"},
		{"tech_research", "TR"},
		{"platform", "PL"},
		{"customization", "CU"},
		{"improvement", "IM"},
		{"others", "OT"},
		{"unknown", "XX"},
	}

	for _, tt := range tests {
		t.Run(tt.category, func(t *testing.T) {
			code := getCategoryCode(tt.category)
			assert.Equal(t, tt.expected, code)
		})
	}
}
