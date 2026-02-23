package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTags(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "valid JSON array",
			input:    `["go", "performance"]`,
			expected: []string{"go", "performance"},
		},
		{
			name:     "empty array",
			input:    "[]",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseTags(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatTags(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "empty slice",
			input:    []string{},
			expected: "[]",
		},
		{
			name:     "valid tags",
			input:    []string{"go", "performance"},
			expected: `["go","performance"]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatTags(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestStringToInt(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		defaultVal int
		expected   int
	}{
		{
			name:       "empty string",
			input:      "",
			defaultVal: 10,
			expected:   10,
		},
		{
			name:       "valid number",
			input:      "20",
			defaultVal: 10,
			expected:   20,
		},
		{
			name:       "invalid number",
			input:      "abc",
			defaultVal: 10,
			expected:   10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringToInt(tt.input, tt.defaultVal)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		item     string
		expected bool
	}{
		{
			name:     "item exists",
			slice:    []string{"a", "b", "c"},
			item:     "b",
			expected: true,
		},
		{
			name:     "item not exists",
			slice:    []string{"a", "b", "c"},
			item:     "d",
			expected: false,
		},
		{
			name:     "case insensitive match",
			slice:    []string{"Go", "Performance"},
			item:     "go",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Contains(tt.slice, tt.item)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractMentions(t *testing.T) {
	service := &ForumService{}

	tests := []struct {
		name     string
		content  string
		expected []string
	}{
		{
			name:     "no mentions",
			content:  "This is a normal post",
			expected: []string{},
		},
		{
			name:     "single mention",
			content:  "Hello @zhangsan",
			expected: []string{"zhangsan"},
		},
		{
			name:     "multiple mentions",
			content:  "@zhangsan @lisi please review",
			expected: []string{"zhangsan", "lisi"},
		},
		{
			name:     "duplicate mentions",
			content:  "@zhangsan @zhangsan",
			expected: []string{"zhangsan"},
		},
		{
			name:     "chinese username",
			content:  "@张三 你好",
			expected: []string{"张三"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.extractMentions(tt.content)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCreateBoardRequestValidation(t *testing.T) {
	req := CreateBoardRequest{
		Name:        "Test Board",
		Description: "Test Description",
		Category:    "tech",
		Icon:        "code",
		SortOrder:   1,
	}

	assert.Equal(t, "Test Board", req.Name)
	assert.Equal(t, "Test Description", req.Description)
	assert.Equal(t, "tech", req.Category)
	assert.Equal(t, "code", req.Icon)
	assert.Equal(t, 1, req.SortOrder)
}

func TestCreatePostRequestValidation(t *testing.T) {
	req := CreatePostRequest{
		Title:   "Test Post",
		Content: "Test Content",
		BoardID: "01J8K2M3N4P5Q6R7S8T9U0V1W",
		Tags:    []string{"test", "go"},
	}

	assert.Equal(t, "Test Post", req.Title)
	assert.Equal(t, "Test Content", req.Content)
	assert.Equal(t, "01J8K2M3N4P5Q6R7S8T9U0V1W", req.BoardID)
	assert.Equal(t, []string{"test", "go"}, req.Tags)
}

func TestCreateReplyRequestValidation(t *testing.T) {
	parentID := "01J8K2M3N4P5Q6R7S8T9U0V1W"
	req := CreateReplyRequest{
		Content:  "Test Reply",
		ParentID: &parentID,
	}

	assert.Equal(t, "Test Reply", req.Content)
	assert.Equal(t, &parentID, req.ParentID)
}

func TestPaginatedResponse(t *testing.T) {
	items := []string{"item1", "item2", "item3"}
	resp := PaginatedResponse{
		Items:    items,
		Total:    10,
		Page:     1,
		PageSize: 3,
	}

	assert.Equal(t, items, resp.Items)
	assert.Equal(t, int64(10), resp.Total)
	assert.Equal(t, 1, resp.Page)
	assert.Equal(t, 3, resp.PageSize)
}

func TestListBoardsQuery(t *testing.T) {
	query := ListBoardsQuery{
		Category: "tech",
		Page:     1,
		PageSize: 20,
	}

	assert.Equal(t, "tech", query.Category)
	assert.Equal(t, 1, query.Page)
	assert.Equal(t, 20, query.PageSize)
}

func TestListPostsQuery(t *testing.T) {
	isPinned := true
	query := ListPostsQuery{
		BoardID:  "board123",
		AuthorID: "user123",
		Search:   "test",
		IsPinned: &isPinned,
		Page:     1,
		PageSize: 20,
	}

	assert.Equal(t, "board123", query.BoardID)
	assert.Equal(t, "user123", query.AuthorID)
	assert.Equal(t, "test", query.Search)
	assert.Equal(t, &isPinned, query.IsPinned)
	assert.Equal(t, 1, query.Page)
	assert.Equal(t, 20, query.PageSize)
}

func TestListRepliesQuery(t *testing.T) {
	parentID := "reply123"
	query := ListRepliesQuery{
		PostID:   "post123",
		ParentID: &parentID,
		Page:     1,
		PageSize: 20,
	}

	assert.Equal(t, "post123", query.PostID)
	assert.Equal(t, &parentID, query.ParentID)
	assert.Equal(t, 1, query.Page)
	assert.Equal(t, 20, query.PageSize)
}
