package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"rdp-platform/rdp-api/models"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// AnalyticsService handles analytics and reporting business logic
type AnalyticsService struct {
	db *gorm.DB
}

// NewAnalyticsService creates a new AnalyticsService
func NewAnalyticsService(db *gorm.DB) *AnalyticsService {
	return &AnalyticsService{db: db}
}

// DashboardOverview represents the overview data for the dashboard
type DashboardOverview struct {
	TotalProjects      int64   `json:"total_projects"`
	ActiveProjects     int64   `json:"active_projects"`
	CompletedProjects  int64   `json:"completed_projects"`
	DelayedProjects    int64   `json:"delayed_projects"`
	TotalUsers         int64   `json:"total_users"`
	ActiveUsers        int64   `json:"active_users"`
	TotalKnowledge     int64   `json:"total_knowledge"`
	TotalProducts      int64   `json:"total_products"`
	AvgProjectProgress float64 `json:"avg_project_progress"`
}

// GetDashboardOverview returns overview statistics for the dashboard
func (s *AnalyticsService) GetDashboardOverview(ctx context.Context) (*DashboardOverview, error) {
	overview := &DashboardOverview{}

	// Count total projects
	if err := s.db.Model(&models.Project{}).Count(&overview.TotalProjects).Error; err != nil {
		return nil, fmt.Errorf("failed to count projects: %w", err)
	}

	// Count active projects (status in active, running)
	if err := s.db.Model(&models.Project{}).
		Where("status IN ?", []string{"active", "running", "in_progress"}).
		Count(&overview.ActiveProjects).Error; err != nil {
		return nil, fmt.Errorf("failed to count active projects: %w", err)
	}

	// Count completed projects
	if err := s.db.Model(&models.Project{}).
		Where("status IN ?", []string{"completed", "done", "closed"}).
		Count(&overview.CompletedProjects).Error; err != nil {
		return nil, fmt.Errorf("failed to count completed projects: %w", err)
	}

	// Count delayed projects (end_date < now and status not completed)
	now := time.Now()
	if err := s.db.Model(&models.Project{}).
		Where("end_date < ? AND status NOT IN ?", now, []string{"completed", "done", "closed"}).
		Count(&overview.DelayedProjects).Error; err != nil {
		return nil, fmt.Errorf("failed to count delayed projects: %w", err)
	}

	// Count total users
	if err := s.db.Model(&models.User{}).Count(&overview.TotalUsers).Error; err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	// Count active users (last_login_at within 30 days)
	thirtyDaysAgo := now.AddDate(0, 0, -30)
	if err := s.db.Model(&models.User{}).
		Where("last_login_at > ?", thirtyDaysAgo).
		Count(&overview.ActiveUsers).Error; err != nil {
		return nil, fmt.Errorf("failed to count active users: %w", err)
	}

	// Count total knowledge
	if err := s.db.Model(&models.Knowledge{}).Count(&overview.TotalKnowledge).Error; err != nil {
		return nil, fmt.Errorf("failed to count knowledge: %w", err)
	}

	// Count total products
	if err := s.db.Model(&models.Product{}).Count(&overview.TotalProducts).Error; err != nil {
		return nil, fmt.Errorf("failed to count products: %w", err)
	}

	// Calculate average project progress
	var avgProgress float64
	if err := s.db.Model(&models.Project{}).
		Select("COALESCE(AVG(progress), 0)").
		Scan(&avgProgress).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate average progress: %w", err)
	}
	overview.AvgProjectProgress = avgProgress

	return overview, nil
}

// ProjectStatistics represents detailed project statistics
type ProjectStatistics struct {
	TotalProjects        int64                `json:"total_projects"`
	StatusDistribution   map[string]int64     `json:"status_distribution"`
	CategoryDistribution map[string]int64     `json:"category_distribution"`
	MonthlyTrend         []MonthlyProjectStat `json:"monthly_trend"`
	TopProjects          []ProjectSummary     `json:"top_projects"`
}

// MonthlyProjectStat represents project statistics for a month
type MonthlyProjectStat struct {
	Month     string `json:"month"`
	Created   int    `json:"created"`
	Completed int    `json:"completed"`
	Active    int    `json:"active"`
}

// ProjectSummary represents a summary of a project for analytics
type ProjectSummary struct {
	ID       string `json:"id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	Progress int    `json:"progress"`
	LeaderID string `json:"leader_id,omitempty"`
}

// GetProjectStatistics returns detailed project statistics
func (s *AnalyticsService) GetProjectStatistics(ctx context.Context, startDate, endDate time.Time) (*ProjectStatistics, error) {
	stats := &ProjectStatistics{
		StatusDistribution:   make(map[string]int64),
		CategoryDistribution: make(map[string]int64),
	}

	// Count total projects in date range
	if err := s.db.Model(&models.Project{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&stats.TotalProjects).Error; err != nil {
		return nil, fmt.Errorf("failed to count projects: %w", err)
	}

	// Get status distribution
	var statusStats []struct {
		Status string
		Count  int64
	}
	if err := s.db.Model(&models.Project{}).
		Select("status, COUNT(*) as count").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Group("status").
		Scan(&statusStats).Error; err != nil {
		return nil, fmt.Errorf("failed to get status distribution: %w", err)
	}
	for _, stat := range statusStats {
		stats.StatusDistribution[stat.Status] = stat.Count
	}

	// Get category distribution
	var categoryStats []struct {
		Category string
		Count    int64
	}
	if err := s.db.Model(&models.Project{}).
		Select("category, COUNT(*) as count").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Group("category").
		Scan(&categoryStats).Error; err != nil {
		return nil, fmt.Errorf("failed to get category distribution: %w", err)
	}
	for _, stat := range categoryStats {
		stats.CategoryDistribution[stat.Category] = stat.Count
	}

	// Get monthly trend
	monthlyStats, err := s.getProjectMonthlyTrend(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}
	stats.MonthlyTrend = monthlyStats

	// Get top projects by progress
	var topProjects []ProjectSummary
	if err := s.db.Model(&models.Project{}).
		Select("id, code, name, status, progress, leader_id").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Order("progress DESC").
		Limit(10).
		Scan(&topProjects).Error; err != nil {
		return nil, fmt.Errorf("failed to get top projects: %w", err)
	}
	stats.TopProjects = topProjects

	return stats, nil
}

// getProjectMonthlyTrend calculates monthly project statistics
func (s *AnalyticsService) getProjectMonthlyTrend(ctx context.Context, startDate, endDate time.Time) ([]MonthlyProjectStat, error) {
	var results []MonthlyProjectStat

	// Generate month list
	current := startDate
	for !current.After(endDate) {
		monthStr := current.Format("2006-01")
		results = append(results, MonthlyProjectStat{Month: monthStr})
		current = current.AddDate(0, 1, 0)
	}

	// Get created projects per month
	var createdStats []struct {
		Month string
		Count int
	}
	if err := s.db.Model(&models.Project{}).
		Select("TO_CHAR(created_at, 'YYYY-MM') as month, COUNT(*) as count").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Group("TO_CHAR(created_at, 'YYYY-MM')").
		Scan(&createdStats).Error; err != nil {
		return nil, fmt.Errorf("failed to get created stats: %w", err)
	}

	// Get completed projects per month
	var completedStats []struct {
		Month string
		Count int
	}
	if err := s.db.Model(&models.Project{}).
		Select("TO_CHAR(actual_end_date, 'YYYY-MM') as month, COUNT(*) as count").
		Where("actual_end_date BETWEEN ? AND ?", startDate, endDate).
		Where("status IN ?", []string{"completed", "done", "closed"}).
		Group("TO_CHAR(actual_end_date, 'YYYY-MM')").
		Scan(&completedStats).Error; err != nil {
		return nil, fmt.Errorf("failed to get completed stats: %w", err)
	}

	// Get active projects per month (simplified: projects active at end of month)
	var activeStats []struct {
		Month string
		Count int
	}
	if err := s.db.Model(&models.Project{}).
		Select("TO_CHAR(created_at, 'YYYY-MM') as month, COUNT(*) as count").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Where("status IN ?", []string{"active", "running", "in_progress"}).
		Group("TO_CHAR(created_at, 'YYYY-MM')").
		Scan(&activeStats).Error; err != nil {
		return nil, fmt.Errorf("failed to get active stats: %w", err)
	}

	// Combine results
	for i := range results {
		for _, stat := range createdStats {
			if stat.Month == results[i].Month {
				results[i].Created = stat.Count
				break
			}
		}
		for _, stat := range completedStats {
			if stat.Month == results[i].Month {
				results[i].Completed = stat.Count
				break
			}
		}
		for _, stat := range activeStats {
			if stat.Month == results[i].Month {
				results[i].Active = stat.Count
				break
			}
		}
	}

	return results, nil
}

// UserStatistics represents user performance statistics
type UserStatistics struct {
	TotalUsers      int64              `json:"total_users"`
	ActiveUsers     int64              `json:"active_users"`
	NewUsers        int64              `json:"new_users"`
	TopContributors []UserContribution `json:"top_contributors"`
	DepartmentStats []DepartmentStat   `json:"department_stats"`
	MonthlyActivity []MonthlyUserStat  `json:"monthly_activity"`
}

// UserContribution represents a user's contribution data
type UserContribution struct {
	UserID         string  `json:"user_id"`
	DisplayName    string  `json:"display_name"`
	AvatarURL      *string `json:"avatar_url,omitempty"`
	ProjectCount   int     `json:"project_count"`
	TaskCount      int     `json:"task_count"`
	KnowledgeCount int     `json:"knowledge_count"`
	Contribution   float64 `json:"contribution"`
}

// DepartmentStat represents statistics for a department
type DepartmentStat struct {
	Department   string `json:"department"`
	UserCount    int    `json:"user_count"`
	ProjectCount int    `json:"project_count"`
}

// MonthlyUserStat represents user statistics for a month
type MonthlyUserStat struct {
	Month       string `json:"month"`
	ActiveUsers int    `json:"active_users"`
	NewUsers    int    `json:"new_users"`
}

// GetUserStatistics returns user performance statistics
func (s *AnalyticsService) GetUserStatistics(ctx context.Context, startDate, endDate time.Time) (*UserStatistics, error) {
	stats := &UserStatistics{}

	// Count total users
	if err := s.db.Model(&models.User{}).Count(&stats.TotalUsers).Error; err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	// Count active users in date range
	if err := s.db.Model(&models.User{}).
		Where("last_login_at BETWEEN ? AND ?", startDate, endDate).
		Count(&stats.ActiveUsers).Error; err != nil {
		return nil, fmt.Errorf("failed to count active users: %w", err)
	}

	// Count new users in date range
	if err := s.db.Model(&models.User{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&stats.NewUsers).Error; err != nil {
		return nil, fmt.Errorf("failed to count new users: %w", err)
	}

	// Get top contributors
	topContributors, err := s.getTopContributors(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}
	stats.TopContributors = topContributors

	// Get department statistics
	deptStats, err := s.getDepartmentStatistics(ctx)
	if err != nil {
		return nil, err
	}
	stats.DepartmentStats = deptStats

	// Get monthly activity
	monthlyStats, err := s.getUserMonthlyActivity(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}
	stats.MonthlyActivity = monthlyStats

	return stats, nil
}

// getTopContributors returns the top contributors based on various metrics
func (s *AnalyticsService) getTopContributors(ctx context.Context, startDate, endDate time.Time) ([]UserContribution, error) {
	var contributors []UserContribution

	// Query to get user contributions
	// This is a simplified version - in production, you might have a dedicated contributions table
	rows, err := s.db.Raw(`
		SELECT 
			u.id as user_id,
			u.display_name,
			u.avatar_url,
			COUNT(DISTINCT pm.project_id) as project_count,
			COUNT(DISTINCT k.id) as knowledge_count
		FROM users u
		LEFT JOIN project_members pm ON pm.user_id = u.id
		LEFT JOIN knowledge k ON k.author_id = u.id AND k.created_at BETWEEN ? AND ?
		GROUP BY u.id, u.display_name, u.avatar_url
		ORDER BY project_count DESC, knowledge_count DESC
		LIMIT 10
	`, startDate, endDate).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to get top contributors: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var c UserContribution
		if err := rows.Scan(&c.UserID, &c.DisplayName, &c.AvatarURL, &c.ProjectCount, &c.KnowledgeCount); err != nil {
			continue
		}
		// Calculate a simple contribution score
		c.Contribution = float64(c.ProjectCount*10 + c.KnowledgeCount*5)
		contributors = append(contributors, c)
	}

	return contributors, nil
}

// getDepartmentStatistics returns statistics by department/team
func (s *AnalyticsService) getDepartmentStatistics(ctx context.Context) ([]DepartmentStat, error) {
	var stats []DepartmentStat

	// Get user count by team
	rows, err := s.db.Raw(`
		SELECT 
			COALESCE(team, '未分配') as department,
			COUNT(*) as user_count
		FROM users
		WHERE is_active = true
		GROUP BY team
		ORDER BY user_count DESC
	`).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to get department stats: %w", err)
	}
	defer rows.Close()

	deptMap := make(map[string]*DepartmentStat)
	for rows.Next() {
		var dept string
		var count int
		if err := rows.Scan(&dept, &count); err != nil {
			continue
		}
		deptMap[dept] = &DepartmentStat{
			Department: dept,
			UserCount:  count,
		}
	}

	// Get project count by team
	projRows, err := s.db.Raw(`
		SELECT 
			COALESCE(team, '未分配') as department,
			COUNT(*) as project_count
		FROM projects
		GROUP BY team
	`).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to get project dept stats: %w", err)
	}
	defer projRows.Close()

	for projRows.Next() {
		var dept string
		var count int
		if err := projRows.Scan(&dept, &count); err != nil {
			continue
		}
		if stat, exists := deptMap[dept]; exists {
			stat.ProjectCount = count
		} else {
			deptMap[dept] = &DepartmentStat{
				Department:   dept,
				ProjectCount: count,
			}
		}
	}

	// Convert map to slice
	for _, stat := range deptMap {
		stats = append(stats, *stat)
	}

	return stats, nil
}

// getUserMonthlyActivity returns monthly user activity statistics
func (s *AnalyticsService) getUserMonthlyActivity(ctx context.Context, startDate, endDate time.Time) ([]MonthlyUserStat, error) {
	var results []MonthlyUserStat

	// Generate month list
	current := startDate
	for !current.After(endDate) {
		monthStr := current.Format("2006-01")
		results = append(results, MonthlyUserStat{Month: monthStr})
		current = current.AddDate(0, 1, 0)
	}

	// Get active users per month
	var activeStats []struct {
		Month string
		Count int
	}
	if err := s.db.Model(&models.User{}).
		Select("TO_CHAR(last_login_at, 'YYYY-MM') as month, COUNT(*) as count").
		Where("last_login_at BETWEEN ? AND ?", startDate, endDate).
		Group("TO_CHAR(last_login_at, 'YYYY-MM')").
		Scan(&activeStats).Error; err != nil {
		return nil, fmt.Errorf("failed to get monthly active users: %w", err)
	}

	// Get new users per month
	var newStats []struct {
		Month string
		Count int
	}
	if err := s.db.Model(&models.User{}).
		Select("TO_CHAR(created_at, 'YYYY-MM') as month, COUNT(*) as count").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Group("TO_CHAR(created_at, 'YYYY-MM')").
		Scan(&newStats).Error; err != nil {
		return nil, fmt.Errorf("failed to get monthly new users: %w", err)
	}

	// Combine results
	for i := range results {
		for _, stat := range activeStats {
			if stat.Month == results[i].Month {
				results[i].ActiveUsers = stat.Count
				break
			}
		}
		for _, stat := range newStats {
			if stat.Month == results[i].Month {
				results[i].NewUsers = stat.Count
				break
			}
		}
	}

	return results, nil
}

// ShelfStatistics represents product shelf statistics
type ShelfStatistics struct {
	TotalProducts     int64          `json:"total_products"`
	PublishedProducts int64          `json:"published_products"`
	TotalTechnologies int64          `json:"total_technologies"`
	CategoryStats     []CategoryStat `json:"category_stats"`
	TopProducts       []ProductUsage `json:"top_products"`
	AdoptionRate      float64        `json:"adoption_rate"`
	ReuseRate         float64        `json:"reuse_rate"`
}

// CategoryStat represents statistics for a category
type CategoryStat struct {
	Category string `json:"category"`
	Count    int    `json:"count"`
	Usage    int    `json:"usage"`
}

// ProductUsage represents usage statistics for a product
type ProductUsage struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	UsageCount  int    `json:"usage_count"`
	CartCount   int    `json:"cart_count"`
}

// GetShelfStatistics returns product shelf statistics
func (s *AnalyticsService) GetShelfStatistics(ctx context.Context, startDate, endDate time.Time) (*ShelfStatistics, error) {
	stats := &ShelfStatistics{}

	// Count total products
	if err := s.db.Model(&models.Product{}).Count(&stats.TotalProducts).Error; err != nil {
		return nil, fmt.Errorf("failed to count products: %w", err)
	}

	// Count published products
	if err := s.db.Model(&models.Product{}).
		Where("is_published = ?", true).
		Count(&stats.PublishedProducts).Error; err != nil {
		return nil, fmt.Errorf("failed to count published products: %w", err)
	}

	// Count total technologies
	if err := s.db.Model(&models.Technology{}).Count(&stats.TotalTechnologies).Error; err != nil {
		return nil, fmt.Errorf("failed to count technologies: %w", err)
	}

	// Get category statistics
	categoryStats, err := s.getProductCategoryStats(ctx)
	if err != nil {
		return nil, err
	}
	stats.CategoryStats = categoryStats

	// Get top products by cart usage
	topProducts, err := s.getTopProducts(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}
	stats.TopProducts = topProducts

	// Calculate adoption rate (published / total)
	if stats.TotalProducts > 0 {
		stats.AdoptionRate = float64(stats.PublishedProducts) / float64(stats.TotalProducts) * 100
	}

	// Calculate reuse rate based on cart items
	var cartItemCount int64
	if err := s.db.Model(&models.CartItem{}).Count(&cartItemCount).Error; err != nil {
		return nil, fmt.Errorf("failed to count cart items: %w", err)
	}
	if stats.TotalProducts > 0 {
		stats.ReuseRate = float64(cartItemCount) / float64(stats.TotalProducts) * 100
	}

	return stats, nil
}

// getProductCategoryStats returns product statistics by category
func (s *AnalyticsService) getProductCategoryStats(ctx context.Context) ([]CategoryStat, error) {
	var stats []CategoryStat

	// For now, we'll use product_line as category since the Product model doesn't have a direct category field
	rows, err := s.db.Raw(`
		SELECT 
			COALESCE(product_line, '未分类') as category,
			COUNT(*) as count
		FROM products
		GROUP BY product_line
		ORDER BY count DESC
	`).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to get category stats: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var stat CategoryStat
		if err := rows.Scan(&stat.Category, &stat.Count); err != nil {
			continue
		}
		stats = append(stats, stat)
	}

	return stats, nil
}

// getTopProducts returns the most used products
func (s *AnalyticsService) getTopProducts(ctx context.Context, startDate, endDate time.Time) ([]ProductUsage, error) {
	var products []ProductUsage

	// Get products with their cart usage
	rows, err := s.db.Raw(`
		SELECT 
			p.id as product_id,
			p.name as product_name,
			COUNT(DISTINCT ci.id) as cart_count
		FROM products p
		LEFT JOIN cart_items ci ON ci.product_id = p.id AND ci.created_at BETWEEN ? AND ?
		GROUP BY p.id, p.name
		ORDER BY cart_count DESC
		LIMIT 10
	`, startDate, endDate).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to get top products: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p ProductUsage
		if err := rows.Scan(&p.ProductID, &p.ProductName, &p.CartCount); err != nil {
			continue
		}
		products = append(products, p)
	}

	return products, nil
}

// KnowledgeStatistics represents knowledge base statistics
type KnowledgeStatistics struct {
	TotalKnowledge int64                   `json:"total_knowledge"`
	PublishedCount int64                   `json:"published_count"`
	DraftCount     int64                   `json:"draft_count"`
	TotalViews     int64                   `json:"total_views"`
	CategoryStats  []KnowledgeCategoryStat `json:"category_stats"`
	TopKnowledge   []KnowledgeSummary      `json:"top_knowledge"`
	TagStats       []TagStat               `json:"tag_stats"`
	MonthlyTrend   []MonthlyKnowledgeStat  `json:"monthly_trend"`
}

// KnowledgeCategoryStat represents knowledge statistics by category
type KnowledgeCategoryStat struct {
	CategoryID   string `json:"category_id"`
	CategoryName string `json:"category_name"`
	Count        int    `json:"count"`
	Views        int    `json:"views"`
}

// KnowledgeSummary represents a summary of knowledge for analytics
type KnowledgeSummary struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	AuthorID  string `json:"author_id"`
	ViewCount int    `json:"view_count"`
}

// TagStat represents statistics for a tag
type TagStat struct {
	TagID   string `json:"tag_id"`
	TagName string `json:"tag_name"`
	Count   int    `json:"count"`
	Color   string `json:"color"`
}

// MonthlyKnowledgeStat represents knowledge statistics for a month
type MonthlyKnowledgeStat struct {
	Month     string `json:"month"`
	Created   int    `json:"created"`
	Published int    `json:"published"`
}

// GetKnowledgeStatistics returns knowledge base statistics
func (s *AnalyticsService) GetKnowledgeStatistics(ctx context.Context, startDate, endDate time.Time) (*KnowledgeStatistics, error) {
	stats := &KnowledgeStatistics{}

	// Count total knowledge
	if err := s.db.Model(&models.Knowledge{}).Count(&stats.TotalKnowledge).Error; err != nil {
		return nil, fmt.Errorf("failed to count knowledge: %w", err)
	}

	// Count by status
	if err := s.db.Model(&models.Knowledge{}).
		Where("status = ?", "published").
		Count(&stats.PublishedCount).Error; err != nil {
		return nil, fmt.Errorf("failed to count published knowledge: %w", err)
	}

	if err := s.db.Model(&models.Knowledge{}).
		Where("status = ?", "draft").
		Count(&stats.DraftCount).Error; err != nil {
		return nil, fmt.Errorf("failed to count draft knowledge: %w", err)
	}

	// Get total views
	var totalViews int64
	if err := s.db.Model(&models.Knowledge{}).
		Select("COALESCE(SUM(view_count), 0)").
		Scan(&totalViews).Error; err != nil {
		return nil, fmt.Errorf("failed to get total views: %w", err)
	}
	stats.TotalViews = totalViews

	// Get category statistics
	categoryStats, err := s.getKnowledgeCategoryStats(ctx)
	if err != nil {
		return nil, err
	}
	stats.CategoryStats = categoryStats

	// Get top knowledge by views
	topKnowledge, err := s.getTopKnowledge(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}
	stats.TopKnowledge = topKnowledge

	// Get tag statistics
	tagStats, err := s.getTagStats(ctx)
	if err != nil {
		return nil, err
	}
	stats.TagStats = tagStats

	// Get monthly trend
	monthlyStats, err := s.getKnowledgeMonthlyTrend(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}
	stats.MonthlyTrend = monthlyStats

	return stats, nil
}

// getKnowledgeCategoryStats returns knowledge statistics by category
func (s *AnalyticsService) getKnowledgeCategoryStats(ctx context.Context) ([]KnowledgeCategoryStat, error) {
	var stats []KnowledgeCategoryStat

	rows, err := s.db.Raw(`
		SELECT 
			c.id as category_id,
			c.name as category_name,
			COUNT(k.id) as count,
			COALESCE(SUM(k.view_count), 0) as views
		FROM categories c
		LEFT JOIN knowledge k ON k.category_id = c.id
		GROUP BY c.id, c.name
		ORDER BY count DESC
	`).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to get knowledge category stats: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var stat KnowledgeCategoryStat
		if err := rows.Scan(&stat.CategoryID, &stat.CategoryName, &stat.Count, &stat.Views); err != nil {
			continue
		}
		stats = append(stats, stat)
	}

	return stats, nil
}

// getTopKnowledge returns the most viewed knowledge entries
func (s *AnalyticsService) getTopKnowledge(ctx context.Context, startDate, endDate time.Time) ([]KnowledgeSummary, error) {
	var knowledge []KnowledgeSummary

	if err := s.db.Model(&models.Knowledge{}).
		Select("id, title, author_id, view_count").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Order("view_count DESC").
		Limit(10).
		Scan(&knowledge).Error; err != nil {
		return nil, fmt.Errorf("failed to get top knowledge: %w", err)
	}

	return knowledge, nil
}

// getTagStats returns statistics for tags
func (s *AnalyticsService) getTagStats(ctx context.Context) ([]TagStat, error) {
	var stats []TagStat

	if err := s.db.Model(&models.Tag{}).
		Select("id as tag_id, name as tag_name, count, color").
		Order("count DESC").
		Limit(20).
		Scan(&stats).Error; err != nil {
		return nil, fmt.Errorf("failed to get tag stats: %w", err)
	}

	return stats, nil
}

// getKnowledgeMonthlyTrend returns monthly knowledge statistics
func (s *AnalyticsService) getKnowledgeMonthlyTrend(ctx context.Context, startDate, endDate time.Time) ([]MonthlyKnowledgeStat, error) {
	var results []MonthlyKnowledgeStat

	// Generate month list
	current := startDate
	for !current.After(endDate) {
		monthStr := current.Format("2006-01")
		results = append(results, MonthlyKnowledgeStat{Month: monthStr})
		current = current.AddDate(0, 1, 0)
	}

	// Get created knowledge per month
	var createdStats []struct {
		Month string
		Count int
	}
	if err := s.db.Model(&models.Knowledge{}).
		Select("TO_CHAR(created_at, 'YYYY-MM') as month, COUNT(*) as count").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Group("TO_CHAR(created_at, 'YYYY-MM')").
		Scan(&createdStats).Error; err != nil {
		return nil, fmt.Errorf("failed to get monthly created knowledge: %w", err)
	}

	// Get published knowledge per month
	var publishedStats []struct {
		Month string
		Count int
	}
	if err := s.db.Model(&models.Knowledge{}).
		Select("TO_CHAR(published_at, 'YYYY-MM') as month, COUNT(*) as count").
		Where("published_at BETWEEN ? AND ?", startDate, endDate).
		Group("TO_CHAR(published_at, 'YYYY-MM')").
		Scan(&publishedStats).Error; err != nil {
		return nil, fmt.Errorf("failed to get monthly published knowledge: %w", err)
	}

	// Combine results
	for i := range results {
		for _, stat := range createdStats {
			if stat.Month == results[i].Month {
				results[i].Created = stat.Count
				break
			}
		}
		for _, stat := range publishedStats {
			if stat.Month == results[i].Month {
				results[i].Published = stat.Count
				break
			}
		}
	}

	return results, nil
}

// DashboardConfig represents a dashboard configuration
type DashboardConfig struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Layout      string    `json:"layout"`
	IsDefault   bool      `json:"is_default"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GetDashboardConfigs returns all dashboard configurations
func (s *AnalyticsService) GetDashboardConfigs(ctx context.Context) ([]DashboardConfig, error) {
	var configs []DashboardConfig

	if err := s.db.Model(&models.AnalyticsDashboard{}).
		Select("id, name, description, layout, is_default, created_by, created_at, updated_at").
		Order("created_at DESC").
		Scan(&configs).Error; err != nil {
		return nil, fmt.Errorf("failed to get dashboard configs: %w", err)
	}

	return configs, nil
}

// GetDashboardConfig returns a dashboard configuration by ID
func (s *AnalyticsService) GetDashboardConfig(ctx context.Context, id string) (*DashboardConfig, error) {
	uid, err := ulid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid dashboard ID")
	}

	var config DashboardConfig
	if err := s.db.Model(&models.AnalyticsDashboard{}).
		Select("id, name, description, layout, is_default, created_by, created_at, updated_at").
		Where("id = ?", uid).
		First(&config).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("dashboard config not found")
		}
		return nil, fmt.Errorf("failed to get dashboard config: %w", err)
	}

	return &config, nil
}

// CreateDashboardConfig creates a new dashboard configuration
func (s *AnalyticsService) CreateDashboardConfig(ctx context.Context, config *DashboardConfig) error {
	dashboard := &models.AnalyticsDashboard{
		Name:        config.Name,
		Description: config.Description,
		Layout:      config.Layout,
		IsDefault:   config.IsDefault,
		CreatedBy:   config.CreatedBy,
	}

	if err := s.db.Create(dashboard).Error; err != nil {
		return fmt.Errorf("failed to create dashboard config: %w", err)
	}

	config.ID = dashboard.ID
	config.CreatedAt = dashboard.CreatedAt
	config.UpdatedAt = dashboard.UpdatedAt

	return nil
}

// UpdateDashboardConfig updates a dashboard configuration
func (s *AnalyticsService) UpdateDashboardConfig(ctx context.Context, id string, config *DashboardConfig) error {
	uid, err := ulid.Parse(id)
	if err != nil {
		return errors.New("invalid dashboard ID")
	}

	result := s.db.Model(
		&models.AnalyticsDashboard{}).
		Where("id = ?", uid).
		Updates(map[string]interface{}{
			"name":        config.Name,
			"description": config.Description,
			"layout":      config.Layout,
			"is_default":  config.IsDefault,
		})

	if result.Error != nil {
		return fmt.Errorf("failed to update dashboard config: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("dashboard config not found")
	}

	return nil
}

// DeleteDashboardConfig deletes a dashboard configuration
func (s *AnalyticsService) DeleteDashboardConfig(ctx context.Context, id string) error {
	uid, err := ulid.Parse(id)
	if err != nil {
		return errors.New("invalid dashboard ID")
	}

	result := s.db.Delete(&models.AnalyticsDashboard{}, "id = ?", uid)
	if result.Error != nil {
		return fmt.Errorf("failed to delete dashboard config: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("dashboard config not found")
	}

	return nil
}

// SetDefaultDashboard sets a dashboard as the default
func (s *AnalyticsService) SetDefaultDashboard(ctx context.Context, id string) error {
	uid, err := ulid.Parse(id)
	if err != nil {
		return errors.New("invalid dashboard ID")
	}

	// Start transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Clear all default flags
	if err := tx.Model(
		&models.AnalyticsDashboard{}).
		Where("is_default = ?", true).
		Update("is_default", false).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to clear default dashboards: %w", err)
	}

	// Set new default
	if err := tx.Model(
		&models.AnalyticsDashboard{}).
		Where("id = ?", uid).
		Update("is_default", true).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to set default dashboard: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GenerateProjectStatsSnapshot generates a snapshot of project statistics for a given date
func (s *AnalyticsService) GenerateProjectStatsSnapshot(ctx context.Context, date time.Time) error {
	// Calculate statistics for the date
	var totalProjects, activeProjects, completedProjects, delayedProjects int64
	var avgProgress float64

	// Count total projects up to the date
	if err := s.db.Model(
		&models.Project{}).
		Where("created_at <= ?", date).
		Count(&totalProjects).Error; err != nil {
		return fmt.Errorf("failed to count total projects: %w", err)
	}

	// Count active projects
	if err := s.db.Model(
		&models.Project{}).
		Where("created_at <= ? AND status IN ?", date, []string{"active", "running", "in_progress"}).
		Count(&activeProjects).Error; err != nil {
		return fmt.Errorf("failed to count active projects: %w", err)
	}

	// Count completed projects up to the date
	if err := s.db.Model(
		&models.Project{}).
		Where("actual_end_date <= ? AND status IN ?", date, []string{"completed", "done", "closed"}).
		Count(&completedProjects).Error; err != nil {
		return fmt.Errorf("failed to count completed projects: %w", err)
	}

	// Count delayed projects
	if err := s.db.Model(
		&models.Project{}).
		Where("created_at <= ? AND end_date < ? AND status NOT IN ?",
			date, date, []string{"completed", "done", "closed"}).
		Count(&delayedProjects).Error; err != nil {
		return fmt.Errorf("failed to count delayed projects: %w", err)
	}

	// Calculate average progress
	if err := s.db.Model(
		&models.Project{}).
		Where("created_at <= ?", date).
		Select("COALESCE(AVG(progress), 0)").
		Scan(&avgProgress).Error; err != nil {
		return fmt.Errorf("failed to calculate average progress: %w", err)
	}

	// Create or update snapshot
	snapshot := &models.ProjectStats{
		Date:              date,
		TotalProjects:     int(totalProjects),
		ActiveProjects:    int(activeProjects),
		CompletedProjects: int(completedProjects),
		DelayedProjects:   int(delayedProjects),
		AvgProgress:       avgProgress,
	}

	// Check if snapshot already exists for this date
	var existing models.ProjectStats
	if err := s.db.Where("date = ?", date).First(&existing).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("failed to check existing snapshot: %w", err)
		}
		// Create new snapshot
		if err := s.db.Create(snapshot).Error; err != nil {
			return fmt.Errorf("failed to create snapshot: %w", err)
		}
	} else {
		// Update existing snapshot
		if err := s.db.Model(
			&existing).Updates(snapshot).Error; err != nil {
			return fmt.Errorf("failed to update snapshot: %w", err)
		}
	}

	return nil
}
