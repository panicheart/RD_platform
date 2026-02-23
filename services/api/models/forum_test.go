package models

import (
	"testing"

	"github.com/oklog/ulid/v2"
)

func TestForumBoardBeforeCreate(t *testing.T) {
	b := &ForumBoard{
		Name:        "Test Board",
		Description: "Test description",
		Category:    "tech",
	}

	err := b.BeforeCreate()
	if err != nil {
		t.Errorf("BeforeCreate() error = %v", err)
	}

	if b.ID == "" {
		t.Error("Expected ID to be generated, got empty string")
	}

	_, err = ulid.Parse(b.ID)
	if err != nil {
		t.Errorf("Generated ID is not a valid ULID: %v", err)
	}
}

func TestForumPostBeforeCreate(t *testing.T) {
	p := &ForumPost{
		BoardID:  "test-board-id",
		Title:    "Test Post",
		Content:  "Test content",
		AuthorID: "test-author-id",
	}

	err := p.BeforeCreate()
	if err != nil {
		t.Errorf("BeforeCreate() error = %v", err)
	}

	if p.ID == "" {
		t.Error("Expected ID to be generated, got empty string")
	}
}

func TestForumReplyBeforeCreate(t *testing.T) {
	r := &ForumReply{
		PostID:   "test-post-id",
		Content:  "Test reply",
		AuthorID: "test-author-id",
	}

	err := r.BeforeCreate()
	if err != nil {
		t.Errorf("BeforeCreate() error = %v", err)
	}

	if r.ID == "" {
		t.Error("Expected ID to be generated, got empty string")
	}
}

func TestForumBoardTableName(t *testing.T) {
	b := ForumBoard{}
	expected := "forum_boards"
	if b.TableName() != expected {
		t.Errorf("TableName() = %v, want %v", b.TableName(), expected)
	}
}

func TestForumPostTableName(t *testing.T) {
	p := ForumPost{}
	expected := "forum_posts"
	if p.TableName() != expected {
		t.Errorf("TableName() = %v, want %v", p.TableName(), expected)
	}
}

func TestForumReplyTableName(t *testing.T) {
	r := ForumReply{}
	expected := "forum_replies"
	if r.TableName() != expected {
		t.Errorf("TableName() = %v, want %v", r.TableName(), expected)
	}
}
