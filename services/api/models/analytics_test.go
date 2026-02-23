package models

import (
	"testing"

	"github.com/oklog/ulid/v2"
)

func TestAnalyticsDashboardBeforeCreate(t *testing.T) {
	d := &AnalyticsDashboard{
		Name:        "Test Dashboard",
		Description: "Test description",
	}

	err := d.BeforeCreate()
	if err != nil {
		t.Errorf("BeforeCreate() error = %v", err)
	}

	if d.ID == "" {
		t.Error("Expected ID to be generated, got empty string")
	}

	_, err = ulid.Parse(d.ID)
	if err != nil {
		t.Errorf("Generated ID is not a valid ULID: %v", err)
	}
}

func TestProjectStatsBeforeCreate(t *testing.T) {
	ps := &ProjectStats{
		TotalProjects:     10,
		ActiveProjects:    5,
		CompletedProjects: 3,
	}

	err := ps.BeforeCreate()
	if err != nil {
		t.Errorf("BeforeCreate() error = %v", err)
	}

	if ps.ID == "" {
		t.Error("Expected ID to be generated, got empty string")
	}
}

func TestUserStatsBeforeCreate(t *testing.T) {
	us := &UserStats{
		UserID:         "test-user-id",
		TasksCompleted: 10,
		Contribution:   85.5,
	}

	err := us.BeforeCreate()
	if err != nil {
		t.Errorf("BeforeCreate() error = %v", err)
	}

	if us.ID == "" {
		t.Error("Expected ID to be generated, got empty string")
	}
}

func TestReportTemplateBeforeCreate(t *testing.T) {
	rt := &ReportTemplate{
		Name:   "Test Report",
		Type:   "project",
		Format: "pdf",
	}

	err := rt.BeforeCreate()
	if err != nil {
		t.Errorf("BeforeCreate() error = %v", err)
	}

	if rt.ID == "" {
		t.Error("Expected ID to be generated, got empty string")
	}
}
