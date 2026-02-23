package models

import (
	"testing"

	"github.com/oklog/ulid/v2"
)

func TestKnowledgeBeforeCreate(t *testing.T) {
	k := &Knowledge{
		Title:   "Test Knowledge",
		Content: "Test content",
	}

	err := k.BeforeCreate()
	if err != nil {
		t.Errorf("BeforeCreate() error = %v", err)
	}

	if k.ID == "" {
		t.Error("Expected ID to be generated, got empty string")
	}

	// Verify it's a valid ULID
	_, err = ulid.Parse(k.ID)
	if err != nil {
		t.Errorf("Generated ID is not a valid ULID: %v", err)
	}
}

func TestCategoryBeforeCreate(t *testing.T) {
	c := &Category{
		Name:        "Test Category",
		Description: "Test description",
		Level:       1,
	}

	err := c.BeforeCreate()
	if err != nil {
		t.Errorf("BeforeCreate() error = %v", err)
	}

	if c.ID == "" {
		t.Error("Expected ID to be generated, got empty string")
	}
}

func TestTagBeforeCreate(t *testing.T) {
	tag := &Tag{
		Name:  "test-tag",
		Color: "#1890ff",
	}

	err := tag.BeforeCreate()
	if err != nil {
		t.Errorf("BeforeCreate() error = %v", err)
	}

	if tag.ID == "" {
		t.Error("Expected ID to be generated, got empty string")
	}
}

func TestKnowledgeTableName(t *testing.T) {
	k := Knowledge{}
	expected := "knowledge"
	if k.TableName() != expected {
		t.Errorf("TableName() = %v, want %v", k.TableName(), expected)
	}
}

func TestCategoryTableName(t *testing.T) {
	c := Category{}
	expected := "categories"
	if c.TableName() != expected {
		t.Errorf("TableName() = %v, want %v", c.TableName(), expected)
	}
}

func TestTagTableName(t *testing.T) {
	tag := Tag{}
	expected := "tags"
	if tag.TableName() != expected {
		t.Errorf("TableName() = %v, want %v", tag.TableName(), expected)
	}
}
