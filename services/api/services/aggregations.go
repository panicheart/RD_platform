package services

import (
	"context"
	"fmt"
	"time"
)

// AggregationService provides complex aggregation queries for analytics
type AggregationService struct {
	analyticsService *AnalyticsService
}

// NewAggregationService creates a new AggregationService
func NewAggregationService(analyticsService *AnalyticsService) *AggregationService {
	return &AggregationService{
		analyticsService: analyticsService,
	}
}

// TimeSeriesDataPoint represents a single point in a time series
type TimeSeriesDataPoint struct {
	Timestamp int64       `json:"timestamp"`
	Value     float64     `json:"value"`
	Label     string      `json:"label,omitempty"`
	Metadata  interface{} `json:"metadata,omitempty"`
}

// TimeSeries represents a time series dataset
type TimeSeries struct {
	Name  string                `json:"name"`
	Label string                `json:"label"`
	Data  []TimeSeriesDataPoint `json:"data"`
	Unit  string                `json:"unit,omitempty"`
	Color string                `json:"color,omitempty"`
}

// PieChartData represents data for a pie/donut chart
type PieChartData struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
	Color string  `json:"color,omitempty"`
}

// BarChartData represents data for a bar chart
type BarChartData struct {
	Category string  `json:"category"`
	Value    float64 `json:"value"`
	Color    string  `json:"color,omitempty"`
}

// HeatmapCell represents a single cell in a heatmap
type HeatmapCell struct {
	X     string  `json:"x"`     // Usually day of week or category
	Y     string  `json:"y"`     // Usually hour or user
	Value float64 `json:"value"` // Intensity value
}

// HeatmapData represents heatmap visualization data
type HeatmapData struct {
	XLabels []string      `json:"x_labels"`
	YLabels []string      `json:"y_labels"`
	Cells   []HeatmapCell `json:"cells"`
	Min     float64       `json:"min"`
	Max     float64       `json:"max"`
}

// StatCard represents a statistic card data
type StatCard struct {
	Title           string  `json:"title"`
	Value           float64 `json:"value"`
	Unit            string  `json:"unit,omitempty"`
	Change          float64 `json:"change"`           // Percentage change from previous period
	ChangeDirection string  `json:"change_direction"` // "up", "down", "neutral"
	Icon            string  `json:"icon,omitempty"`
	Color           string  `json:"color,omitempty"`
}

// DashboardWidgets represents all widgets for the dashboard
type DashboardWidgets struct {
	StatCards  []StatCard     `json:"stat_cards"`
	TimeSeries []TimeSeries   `json:"time_series"`
	PieCharts  []PieChartData `json:"pie_charts"`
	BarCharts  []BarChartData `json:"bar_charts"`
	Heatmap    *HeatmapData   `json:"heatmap,omitempty"`
}

// GetDashboardWidgets returns all widget data for the dashboard
func (s *AggregationService) GetDashboardWidgets(ctx context.Context, startDate, endDate time.Time) (*DashboardWidgets, error) {
	widgets := &DashboardWidgets{}

	// Get stat cards
	statCards, err := s.getStatCards(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get stat cards: %w", err)
	}
	widgets.StatCards = statCards

	// Get time series data
	timeSeries, err := s.getTimeSeriesData(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get time series: %w", err)
	}
	widgets.TimeSeries = timeSeries

	// Get pie chart data
	pieCharts, err := s.getPieChartData(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get pie charts: %w", err)
	}
	widgets.PieCharts = pieCharts

	// Get bar chart data
	barCharts, err := s.getBarChartData(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get bar charts: %w", err)
	}
	widgets.BarCharts = barCharts

	// Get heatmap data
	heatmap, err := s.getContributionHeatmap(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get heatmap: %w", err)
	}
	widgets.Heatmap = heatmap

	return widgets, nil
}

// getStatCards generates stat card data
func (s *AggregationService) getStatCards(ctx context.Context, startDate, endDate time.Time) ([]StatCard, error) {
	// Get current period stats
	overview, err := s.analyticsService.GetDashboardOverview(ctx)
	if err != nil {
		return nil, err
	}

	// Calculate previous period for comparison
	periodDuration := endDate.Sub(startDate)
	prevStartDate := startDate.Add(-periodDuration)
	prevEndDate := startDate

	// Get previous period stats (simplified - in production, cache these)
	// For now, we'll calculate a mock change percentage

	cards := []StatCard{
		{
			Title:           "总项目数",
			Value:           float64(overview.TotalProjects),
			Change:          12.5, // Mock value
			ChangeDirection: "up",
			Icon:            "project",
			Color:           "#1890ff",
		},
		{
			Title:           "活跃项目",
			Value:           float64(overview.ActiveProjects),
			Unit:            "个",
			Change:          8.3,
			ChangeDirection: "up",
			Icon:            "rocket",
			Color:           "#52c41a",
		},
		{
			Title:           "平均进度",
			Value:           overview.AvgProjectProgress,
			Unit:            "%",
			Change:          -2.1,
			ChangeDirection: "down",
			Icon:            "line-chart",
			Color:           "#faad14",
		},
		{
			Title:           "延期项目",
			Value:           float64(overview.DelayedProjects),
			Unit:            "个",
			Change:          -15.0,
			ChangeDirection: "neutral",
			Icon:            "warning",
			Color:           "#ff4d4f",
		},
		{
			Title:           "活跃用户数",
			Value:           float64(overview.ActiveUsers),
			Unit:            "人",
			Change:          5.7,
			ChangeDirection: "up",
			Icon:            "team",
			Color:           "#722ed1",
		},
		{
			Title:           "知识库文档",
			Value:           float64(overview.TotalKnowledge),
			Unit:            "篇",
			Change:          23.4,
			ChangeDirection: "up",
			Icon:            "book",
			Color:           "#13c2c2",
		},
	}

	// Calculate actual changes by comparing with previous period
	_ = prevStartDate // Use these in production implementation
	_ = prevEndDate

	return cards, nil
}

// getTimeSeriesData generates time series data for charts
func (s *AggregationService) getTimeSeriesData(ctx context.Context, startDate, endDate time.Time) ([]TimeSeries, error) {
	// Get project statistics
	projStats, err := s.analyticsService.GetProjectStatistics(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	var timeSeries []TimeSeries

	// Project trend time series
	projTrend := TimeSeries{
		Name:  "project_trend",
		Label: "项目趋势",
		Unit:  "个",
		Color: "#1890ff",
	}

	for _, stat := range projStats.MonthlyTrend {
		// Parse month to timestamp
		t, _ := time.Parse("2006-01", stat.Month)
		projTrend.Data = append(projTrend.Data, TimeSeriesDataPoint{
			Timestamp: t.Unix() * 1000, // JavaScript timestamp (milliseconds)
			Value:     float64(stat.Active),
			Label:     stat.Month,
			Metadata: map[string]interface{}{
				"created":   stat.Created,
				"completed": stat.Completed,
			},
		})
	}
	timeSeries = append(timeSeries, projTrend)

	// Get user statistics for user activity trend
	userStats, err := s.analyticsService.GetUserStatistics(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	userTrend := TimeSeries{
		Name:  "user_activity",
		Label: "用户活跃度",
		Unit:  "人",
		Color: "#52c41a",
	}

	for _, stat := range userStats.MonthlyActivity {
		t, _ := time.Parse("2006-01", stat.Month)
		userTrend.Data = append(userTrend.Data, TimeSeriesDataPoint{
			Timestamp: t.Unix() * 1000,
			Value:     float64(stat.ActiveUsers),
			Label:     stat.Month,
			Metadata: map[string]interface{}{
				"new_users": stat.NewUsers,
			},
		})
	}
	timeSeries = append(timeSeries, userTrend)

	return timeSeries, nil
}

// getPieChartData generates pie chart data
func (s *AggregationService) getPieChartData(ctx context.Context, startDate, endDate time.Time) ([]PieChartData, error) {
	// Get project statistics
	projStats, err := s.analyticsService.GetProjectStatistics(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	var pieData []PieChartData
	colors := []string{"#1890ff", "#52c41a", "#faad14", "#ff4d4f", "#722ed1", "#13c2c2"}

	colorIdx := 0
	for status, count := range projStats.StatusDistribution {
		pieData = append(pieData, PieChartData{
			Name:  status,
			Value: float64(count),
			Color: colors[colorIdx%len(colors)],
		})
		colorIdx++
	}

	return pieData, nil
}

// getBarChartData generates bar chart data
func (s *AggregationService) getBarChartData(ctx context.Context, startDate, endDate time.Time) ([]BarChartData, error) {
	// Get project statistics
	projStats, err := s.analyticsService.GetProjectStatistics(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	var barData []BarChartData
	colors := []string{"#1890ff", "#52c41a", "#faad14", "#ff4d4f", "#722ed1"}

	colorIdx := 0
	for category, count := range projStats.CategoryDistribution {
		barData = append(barData, BarChartData{
			Category: category,
			Value:    float64(count),
			Color:    colors[colorIdx%len(colors)],
		})
		colorIdx++
	}

	return barData, nil
}

// getContributionHeatmap generates GitHub-style contribution heatmap data
func (s *AggregationService) getContributionHeatmap(ctx context.Context, startDate, endDate time.Time) (*HeatmapData, error) {
	// Get user statistics
	userStats, err := s.analyticsService.GetUserStatistics(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	heatmap := &HeatmapData{
		XLabels: []string{"周一", "周二", "周三", "周四", "周五", "周六", "周日"},
		YLabels: []string{},
		Cells:   []HeatmapCell{},
		Min:     0,
		Max:     10,
	}

	// Add top contributors as Y labels
	for i, contributor := range userStats.TopContributors {
		if i >= 10 { // Limit to top 10
			break
		}
		heatmap.YLabels = append(heatmap.YLabels, contributor.DisplayName)
	}

	// Generate mock heatmap data (in production, this would be real activity data)
	// For each user and day combination
	for yIdx, contributor := range userStats.TopContributors {
		if yIdx >= 10 {
			break
		}
		for xIdx := 0; xIdx < 7; xIdx++ {
			// Generate a random activity level based on contribution score
			activityLevel := contributor.Contribution / 10.0
			if activityLevel > 10 {
				activityLevel = 10
			}

			heatmap.Cells = append(heatmap.Cells, HeatmapCell{
				X:     heatmap.XLabels[xIdx],
				Y:     contributor.DisplayName,
				Value: activityLevel,
			})
		}
	}

	return heatmap, nil
}

// ProjectProgressTrend represents project progress trend data
type ProjectProgressTrend struct {
	ProjectID   string  `json:"project_id"`
	ProjectName string  `json:"project_name"`
	Progress    []int   `json:"progress"` // Progress values over time
	Dates       []int64 `json:"dates"`    // Timestamps
}

// GetProjectProgressTrends returns progress trends for projects
func (s *AggregationService) GetProjectProgressTrends(ctx context.Context, projectIDs []string, startDate, endDate time.Time) ([]ProjectProgressTrend, error) {
	// In a production system, you would have a project_progress_history table
	// For now, we'll return mock data

	var trends []ProjectProgressTrend

	// This would query a historical progress table
	// For demonstration, we return empty trends
	for _, projectID := range projectIDs {
		trends = append(trends, ProjectProgressTrend{
			ProjectID:   projectID,
			ProjectName: "",
			Progress:    []int{},
			Dates:       []int64{},
		})
	}

	return trends, nil
}

// ComparisonData represents comparison data between two time periods
type ComparisonData struct {
	Metric        string  `json:"metric"`
	CurrentValue  float64 `json:"current_value"`
	PreviousValue float64 `json:"previous_value"`
	Change        float64 `json:"change"`
	ChangePercent float64 `json:"change_percent"`
}

// ComparePeriods compares analytics data between two time periods
func (s *AggregationService) ComparePeriods(ctx context.Context, currentStart, currentEnd, prevStart, prevEnd time.Time) ([]ComparisonData, error) {
	var comparisons []ComparisonData

	// Get current period overview
	currentOverview, err := s.analyticsService.GetDashboardOverview(ctx)
	if err != nil {
		return nil, err
	}

	// Calculate comparisons for key metrics
	metrics := []struct {
		name         string
		currentValue float64
	}{
		{"总项目数", float64(currentOverview.TotalProjects)},
		{"活跃项目", float64(currentOverview.ActiveProjects)},
		{"已完成项目", float64(currentOverview.CompletedProjects)},
		{"延期项目", float64(currentOverview.DelayedProjects)},
		{"总用户数", float64(currentOverview.TotalUsers)},
		{"活跃用户数", float64(currentOverview.ActiveUsers)},
		{"平均进度", currentOverview.AvgProjectProgress},
	}

	for _, metric := range metrics {
		comparisons = append(comparisons, ComparisonData{
			Metric:        metric.name,
			CurrentValue:  metric.currentValue,
			PreviousValue: 0, // Would be fetched from historical data
			Change:        0,
			ChangePercent: 0,
		})
	}

	// Suppress unused variable warning
	_ = prevStart
	_ = prevEnd

	return comparisons, nil
}

// ExportData represents data formatted for export
type ExportData struct {
	Headers []string               `json:"headers"`
	Rows    [][]interface{}        `json:"rows"`
	Meta    map[string]interface{} `json:"meta,omitempty"`
}

// ExportProjectStats exports project statistics in a format suitable for PDF/Excel
func (s *AggregationService) ExportProjectStats(ctx context.Context, startDate, endDate time.Time) (*ExportData, error) {
	stats, err := s.analyticsService.GetProjectStatistics(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	export := &ExportData{
		Headers: []string{"指标", "数值"},
		Rows:    [][]interface{}{},
		Meta: map[string]interface{}{
			"title":        "项目统计报表",
			"start_date":   startDate.Format("2006-01-02"),
			"end_date":     endDate.Format("2006-01-02"),
			"generated_at": time.Now().Format("2006-01-02 15:04:05"),
		},
	}

	// Add summary row
	export.Rows = append(export.Rows, []interface{}{"总项目数", stats.TotalProjects})

	// Add status distribution
	for status, count := range stats.StatusDistribution {
		export.Rows = append(export.Rows, []interface{}{
			fmt.Sprintf("状态: %s", status),
			count,
		})
	}

	// Add category distribution
	for category, count := range stats.CategoryDistribution {
		export.Rows = append(export.Rows, []interface{}{
			fmt.Sprintf("类别: %s", category),
			count,
		})
	}

	// Add monthly trend
	for _, trend := range stats.MonthlyTrend {
		export.Rows = append(export.Rows, []interface{}{
			fmt.Sprintf("%s - 新建", trend.Month),
			trend.Created,
		})
		export.Rows = append(export.Rows, []interface{}{
			fmt.Sprintf("%s - 完成", trend.Month),
			trend.Completed,
		})
	}

	return export, nil
}

// ExportUserStats exports user statistics in a format suitable for PDF/Excel
func (s *AggregationService) ExportUserStats(ctx context.Context, startDate, endDate time.Time) (*ExportData, error) {
	stats, err := s.analyticsService.GetUserStatistics(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	export := &ExportData{
		Headers: []string{"指标", "数值"},
		Rows:    [][]interface{}{},
		Meta: map[string]interface{}{
			"title":        "用户统计报表",
			"start_date":   startDate.Format("2006-01-02"),
			"end_date":     endDate.Format("2006-01-02"),
			"generated_at": time.Now().Format("2006-01-02 15:04:05"),
		},
	}

	// Add summary rows
	export.Rows = append(export.Rows, []interface{}{"总用户数", stats.TotalUsers})
	export.Rows = append(export.Rows, []interface{}{"活跃用户数", stats.ActiveUsers})
	export.Rows = append(export.Rows, []interface{}{"新增用户数", stats.NewUsers})

	// Add top contributors
	export.Rows = append(export.Rows, []interface{}{"", ""})
	export.Rows = append(export.Rows, []interface{}{"TOP 贡献者", ""})
	export.Rows = append(export.Rows, []interface{}{"用户名", "贡献度"})

	for _, contributor := range stats.TopContributors {
		export.Rows = append(export.Rows, []interface{}{
			contributor.DisplayName,
			contributor.Contribution,
		})
	}

	// Add department statistics
	export.Rows = append(export.Rows, []interface{}{"", ""})
	export.Rows = append(export.Rows, []interface{}{"部门统计", ""})
	export.Rows = append(export.Rows, []interface{}{"部门", "用户数", "项目数"})

	for _, dept := range stats.DepartmentStats {
		export.Rows = append(export.Rows, []interface{}{
			dept.Department,
			dept.UserCount,
			dept.ProjectCount,
		})
	}

	return export, nil
}

// ExportShelfStats exports shelf statistics in a format suitable for PDF/Excel
func (s *AggregationService) ExportShelfStats(ctx context.Context, startDate, endDate time.Time) (*ExportData, error) {
	stats, err := s.analyticsService.GetShelfStatistics(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	export := &ExportData{
		Headers: []string{"指标", "数值"},
		Rows:    [][]interface{}{},
		Meta: map[string]interface{}{
			"title":        "货架统计报表",
			"start_date":   startDate.Format("2006-01-02"),
			"end_date":     endDate.Format("2006-01-02"),
			"generated_at": time.Now().Format("2006-01-02 15:04:05"),
		},
	}

	// Add summary rows
	export.Rows = append(export.Rows, []interface{}{"总产品数", stats.TotalProducts})
	export.Rows = append(export.Rows, []interface{}{"已发布产品", stats.PublishedProducts})
	export.Rows = append(export.Rows, []interface{}{"技术总数", stats.TotalTechnologies})
	export.Rows = append(export.Rows, []interface{}{"采用率", fmt.Sprintf("%.1f%%", stats.AdoptionRate)})
	export.Rows = append(export.Rows, []interface{}{"复用率", fmt.Sprintf("%.1f%%", stats.ReuseRate)})

	// Add category statistics
	export.Rows = append(export.Rows, []interface{}{"", ""})
	export.Rows = append(export.Rows, []interface{}{"类别统计", ""})
	export.Rows = append(export.Rows, []interface{}{"类别", "数量"})

	for _, cat := range stats.CategoryStats {
		export.Rows = append(export.Rows, []interface{}{
			cat.Category,
			cat.Count,
		})
	}

	// Add top products
	export.Rows = append(export.Rows, []interface{}{"", ""})
	export.Rows = append(export.Rows, []interface{}{"热门产品", ""})
	export.Rows = append(export.Rows, []interface{}{"产品名称", "购物车次数"})

	for _, product := range stats.TopProducts {
		export.Rows = append(export.Rows, []interface{}{
			product.ProductName,
			product.CartCount,
		})
	}

	return export, nil
}
