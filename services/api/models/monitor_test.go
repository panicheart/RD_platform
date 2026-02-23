package models

import (
	"testing"

	"github.com/oklog/ulid/v2"
)

func TestSystemMetricBeforeCreate(t *testing.T) {
	sm := &SystemMetric{
		CPUUsage:    50.5,
		MemoryUsage: 75.0,
	}

	err := sm.BeforeCreate()
	if err != nil {
		t.Errorf("BeforeCreate() error = %v", err)
	}

	if sm.ID == "" {
		t.Error("Expected ID to be generated, got empty string")
	}

	_, err = ulid.Parse(sm.ID)
	if err != nil {
		t.Errorf("Generated ID is not a valid ULID: %v", err)
	}
}

func TestAPIMetricBeforeCreate(t *testing.T) {
	am := &APIMetric{
		Endpoint:   "/api/v1/test",
		Method:     "GET",
		Duration:   100,
		StatusCode: 200,
	}

	err := am.BeforeCreate()
	if err != nil {
		t.Errorf("BeforeCreate() error = %v", err)
	}

	if am.ID == "" {
		t.Error("Expected ID to be generated, got empty string")
	}
}

func TestLogEntryBeforeCreate(t *testing.T) {
	le := &LogEntry{
		Level:   "INFO",
		Message: "Test log message",
		Source:  "test-service",
	}

	err := le.BeforeCreate()
	if err != nil {
		t.Errorf("BeforeCreate() error = %v", err)
	}

	if le.ID == "" {
		t.Error("Expected ID to be generated, got empty string")
	}
}

func TestAlertRuleBeforeCreate(t *testing.T) {
	ar := &AlertRule{
		Name:      "Test Alert",
		Metric:    "cpu_usage",
		Condition: ">",
		Threshold: 80.0,
	}

	err := ar.BeforeCreate()
	if err != nil {
		t.Errorf("BeforeCreate() error = %v", err)
	}

	if ar.ID == "" {
		t.Error("Expected ID to be generated, got empty string")
	}
}

func TestAlertHistoryBeforeCreate(t *testing.T) {
	ah := &AlertHistory{
		RuleID:    "test-rule-id",
		RuleName:  "Test Rule",
		Severity:  "warning",
		Status:    "firing",
		Value:     85.0,
		Threshold: 80.0,
	}

	err := ah.BeforeCreate()
	if err != nil {
		t.Errorf("BeforeCreate() error = %v", err)
	}

	if ah.ID == "" {
		t.Error("Expected ID to be generated, got empty string")
	}
}
