package services

import (
	"context"
	"testing"
	"time"

	"rdp-platform/rdp-api/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type AnalyticsServiceTestSuite struct {
	suite.Suite
	db               *gorm.DB
	analyticsService *AnalyticsService
}

func (s *AnalyticsServiceTestSuite) SetupTest() {
	var err error
	s.db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	s.Require().NoError(err)

	// Migrate tables
	s.db.AutoMigrate(
		&models.Project{},
		&models.User{},
		&models.Knowledge{},
		&models.Product{},
		&models.Technology{},
		&models.Category{},
		&models.Tag{},
		&models.AnalyticsDashboard{},
		&models.ProjectStats{},
	)

	s.analyticsService = NewAnalyticsService(s.db)
}

func (s *AnalyticsServiceTestSuite) TearDownTest() {
	sqlDB, err := s.db.DB()
	if err == nil {
		sqlDB.Close()
	}
}

func TestAnalyticsServiceSuite(t *testing.T) {
	suite.Run(t, new(AnalyticsServiceTestSuite))
}

func (s *AnalyticsServiceTestSuite) TestGetDashboardOverview() {
	ctx := context.Background()

	// Create test data
	now := time.Now()
	s.db.Create(&models.Project{
		Name:   "Test Project 1",
		Status: "active",
		Code:   "RDP-TEST-001",
	})
	s.db.Create(&models.Project{
		Name:   "Test Project 2",
		Status: "completed",
		Code:   "RDP-TEST-002",
	})
	s.db.Create(&models.User{
		Username:    "testuser",
		DisplayName: "Test User",
		LastLoginAt: &now,
	})
	s.db.Create(&models.Knowledge{
		Title:  "Test Knowledge",
		Status: "published",
	})
	s.db.Create(&models.Product{
		Name:     "Test Product",
		Category: "Test Category",
	})

	// Test GetDashboardOverview
	overview, err := s.analyticsService.GetDashboardOverview(ctx)
	s.Require().NoError(err)
	s.NotNil(overview)

	// Verify counts
	s.Equal(int64(2), overview.TotalProjects)
	s.Equal(int64(1), overview.ActiveProjects)
	s.Equal(int64(1), overview.CompletedProjects)
	s.Equal(int64(1), overview.TotalUsers)
	s.Equal(int64(1), overview.TotalKnowledge)
	s.Equal(int64(1), overview.TotalProducts)
}

func (s *AnalyticsServiceTestSuite) TestGetProjectStatistics() {
	ctx := context.Background()
	startDate := time.Now().AddDate(0, -1, 0)
	endDate := time.Now().AddDate(0, 1, 0)

	// Create test projects
	s.db.Create(&models.Project{
		Name:     "Active Project",
		Status:   "active",
		Code:     "RDP-TEST-001",
		Category: "R&D",
	})
	s.db.Create(&models.Project{
		Name:     "Completed Project",
		Status:   "completed",
		Code:     "RDP-TEST-002",
		Category: "Product",
	})

	stats, err := s.analyticsService.GetProjectStatistics(ctx, startDate, endDate)
	s.Require().NoError(err)
	s.NotNil(stats)

	s.Equal(int64(2), stats.TotalProjects)
	s.NotNil(stats.StatusDistribution)
	s.NotNil(stats.CategoryDistribution)
}

func (s *AnalyticsServiceTestSuite) TestGetUserStatistics() {
	ctx := context.Background()
	startDate := time.Now().AddDate(0, -1, 0)
	endDate := time.Now().AddDate(0, 1, 0)

	// Create test users
	now := time.Now()
	teamA := "Team A"
	s.db.Create(&models.User{
		Username:    "user1",
		DisplayName: "User One",
		LastLoginAt: &now,
		Team:        &teamA,
	})
	s.db.Create(&models.User{
		Username:    "user2",
		DisplayName: "User Two",
		LastLoginAt: &now,
		Team:        &teamA,
	})

	stats, err := s.analyticsService.GetUserStatistics(ctx, startDate, endDate)
	s.Require().NoError(err)
	s.NotNil(stats)

	s.Equal(int64(2), stats.TotalUsers)
}

func (s *AnalyticsServiceTestSuite) TestGetShelfStatistics() {
	ctx := context.Background()
	startDate := time.Now().AddDate(0, -1, 0)
	endDate := time.Now().AddDate(0, 1, 0)

	// Create test products
	s.db.Create(&models.Product{
		Name:        "Product 1",
		IsPublished: true,
		Category:    "Line A",
	})
	s.db.Create(&models.Product{
		Name:        "Product 2",
		IsPublished: false,
		Category:    "Line B",
	})
	s.db.Create(&models.Technology{
		Name:     "Technology 1",
		Category: "Test Category",
	})

	stats, err := s.analyticsService.GetShelfStatistics(ctx, startDate, endDate)
	s.Require().NoError(err)
	s.NotNil(stats)

	s.Equal(int64(2), stats.TotalProducts)
	s.Equal(int64(1), stats.PublishedProducts)
	s.Equal(int64(1), stats.TotalTechnologies)
}

func (s *AnalyticsServiceTestSuite) TestGetKnowledgeStatistics() {
	ctx := context.Background()
	startDate := time.Now().AddDate(0, -1, 0)
	endDate := time.Now().AddDate(0, 1, 0)

	// Create test category and knowledge
	s.db.Create(&models.Category{
		Name: "Test Category",
	})
	s.db.Create(&models.Knowledge{
		Title:     "Published Knowledge",
		Status:    "published",
		ViewCount: 100,
	})
	s.db.Create(&models.Knowledge{
		Title:  "Draft Knowledge",
		Status: "draft",
	})
	s.db.Create(&models.Tag{
		Name:  "Test Tag",
		Count: 5,
	})

	stats, err := s.analyticsService.GetKnowledgeStatistics(ctx, startDate, endDate)
	s.Require().NoError(err)
	s.NotNil(stats)

	s.Equal(int64(2), stats.TotalKnowledge)
	s.Equal(int64(1), stats.PublishedCount)
	s.Equal(int64(1), stats.DraftCount)
	s.Equal(int64(100), stats.TotalViews)
}

func (s *AnalyticsServiceTestSuite) TestDashboardConfigCRUD() {
	ctx := context.Background()

	// Test Create
	config := &DashboardConfig{
		Name:        "Test Dashboard",
		Description: "Test Description",
		Layout:      "{}",
		IsDefault:   false,
		CreatedBy:   "user1",
	}

	err := s.analyticsService.CreateDashboardConfig(ctx, config)
	s.Require().NoError(err)
	s.NotEmpty(config.ID)

	// Test Get
	retrieved, err := s.analyticsService.GetDashboardConfig(ctx, config.ID)
	s.Require().NoError(err)
	s.Equal(config.Name, retrieved.Name)

	// Test List
	configs, err := s.analyticsService.GetDashboardConfigs(ctx)
	s.Require().NoError(err)
	s.Len(configs, 1)

	// Test Update
	config.Name = "Updated Dashboard"
	err = s.analyticsService.UpdateDashboardConfig(ctx, config.ID, config)
	s.Require().NoError(err)

	retrieved, err = s.analyticsService.GetDashboardConfig(ctx, config.ID)
	s.Require().NoError(err)
	s.Equal("Updated Dashboard", retrieved.Name)

	// Test Delete
	err = s.analyticsService.DeleteDashboardConfig(ctx, config.ID)
	s.Require().NoError(err)

	_, err = s.analyticsService.GetDashboardConfig(ctx, config.ID)
	s.Error(err)
}

func (s *AnalyticsServiceTestSuite) TestSetDefaultDashboard() {
	ctx := context.Background()

	// Create two dashboards
	config1 := &DashboardConfig{
		Name:      "Dashboard 1",
		Layout:    "{}",
		IsDefault: false,
	}
	config2 := &DashboardConfig{
		Name:      "Dashboard 2",
		Layout:    "{}",
		IsDefault: false,
	}

	s.Require().NoError(s.analyticsService.CreateDashboardConfig(ctx, config1))
	s.Require().NoError(s.analyticsService.CreateDashboardConfig(ctx, config2))

	// Set first as default
	err := s.analyticsService.SetDefaultDashboard(ctx, config1.ID)
	s.Require().NoError(err)

	// Verify
	retrieved, err := s.analyticsService.GetDashboardConfig(ctx, config1.ID)
	s.Require().NoError(err)
	s.True(retrieved.IsDefault)

	// Set second as default
	err = s.analyticsService.SetDefaultDashboard(ctx, config2.ID)
	s.Require().NoError(err)

	// Verify first is no longer default
	retrieved, err = s.analyticsService.GetDashboardConfig(ctx, config1.ID)
	s.Require().NoError(err)
	s.False(retrieved.IsDefault)

	// Verify second is default
	retrieved, err = s.analyticsService.GetDashboardConfig(ctx, config2.ID)
	s.Require().NoError(err)
	s.True(retrieved.IsDefault)
}

func (s *AnalyticsServiceTestSuite) TestGenerateProjectStatsSnapshot() {
	ctx := context.Background()
	date := time.Now().Truncate(24 * time.Hour)

	// Create test projects
	s.db.Create(&models.Project{
		Name:   "Project 1",
		Status: "active",
		Code:   "RDP-001",
	})
	s.db.Create(&models.Project{
		Name:   "Project 2",
		Status: "completed",
		Code:   "RDP-002",
	})

	// Generate snapshot
	err := s.analyticsService.GenerateProjectStatsSnapshot(ctx, date)
	s.Require().NoError(err)

	// Verify snapshot was created
	var snapshot models.ProjectStats
	result := s.db.Where("date = ?", date).First(&snapshot)
	s.NoError(result.Error)
	s.Equal(2, snapshot.TotalProjects)
}

func TestAnalyticsService(t *testing.T) {
	// Basic unit tests

	t.Run("NewAnalyticsService", func(t *testing.T) {
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)
		defer func() {
			sqlDB, _ := db.DB()
			sqlDB.Close()
		}()

		service := NewAnalyticsService(db)
		assert.NotNil(t, service)
	})
}
