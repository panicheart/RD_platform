package services

import (
	"testing"
	"time"
)

func TestParseFrontmatter(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantMeta map[string]string
		wantBody string
	}{
		{
			name: "simple frontmatter",
			content: `---
title: Test Document
author: user123
tags: tag1, tag2
---
This is the body content.`,
			wantMeta: map[string]string{
				"title":  "Test Document",
				"author": "user123",
				"tags":   "tag1, tag2",
			},
			wantBody: "This is the body content.",
		},
		{
			name:     "no frontmatter",
			content:  "Just content without frontmatter",
			wantMeta: map[string]string{},
			wantBody: "Just content without frontmatter",
		},
		{
			name: "empty frontmatter",
			content: `---
---
Body only`,
			wantMeta: map[string]string{},
			wantBody: "Body only",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta, body := parseFrontmatter(tt.content)
			
			// Check metadata
			for key, want := range tt.wantMeta {
				if got := meta[key]; got != want {
					t.Errorf("metadata[%q] = %q, want %q", key, got, want)
				}
			}
			
			// Check body
			if body != tt.wantBody {
				t.Errorf("body = %q, want %q", body, tt.wantBody)
			}
		})
	}
}

func TestExtractTags(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []string
	}{
		{
			name:    "single tag",
			content: "This is a document with #tag1",
			want:    []string{"tag1"},
		},
		{
			name:    "multiple tags",
			content: "Document with #tag1 and #tag2 #tag3",
			want:    []string{"tag1", "tag2", "tag3"},
		},
		{
			name:    "duplicate tags",
			content: "Document with #tag1 and #tag1 again",
			want:    []string{"tag1"},
		},
		{
			name:    "chinese tags",
			content: "文档包含 #中文标签 和 #微波技术",
			want:    []string{"中文标签", "微波技术"},
		},
		{
			name:    "no tags",
			content: "Document without any tags",
			want:    []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractTags(tt.content)
			if len(got) != len(tt.want) {
				t.Errorf("extractTags() = %v, want %v", got, tt.want)
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("extractTags()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestMergeTags(t *testing.T) {
	tests := []struct {
		name string
		a    []string
		b    []string
		want []string
	}{
		{
			name: "no duplicates",
			a:    []string{"tag1", "tag2"},
			b:    []string{"tag3", "tag4"},
			want: []string{"tag1", "tag2", "tag3", "tag4"},
		},
		{
			name: "with duplicates",
			a:    []string{"tag1", "tag2"},
			b:    []string{"tag2", "tag3"},
			want: []string{"tag1", "tag2", "tag3"},
		},
		{
			name: "empty a",
			a:    []string{},
			b:    []string{"tag1", "tag2"},
			want: []string{"tag1", "tag2"},
		},
		{
			name: "empty b",
			a:    []string{"tag1", "tag2"},
			b:    []string{},
			want: []string{"tag1", "tag2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mergeTags(tt.a, tt.b)
			if len(got) != len(tt.want) {
				t.Errorf("mergeTags() = %v, want %v", got, tt.want)
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("mergeTags()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "normal filename",
			in:   "normal filename",
			want: "normal filename",
		},
		{
			name: "with special chars",
			in:   "file:name/with|special?chars",
			want: "file_name_with_special_chars",
		},
		{
			name: "with spaces",
			in:   "  file with spaces  ",
			want: "file with spaces",
		},
		{
			name: "very long",
			in:   "a very long filename that exceeds the maximum allowed length of two hundred characters for obsidian vault file names and should be truncated",
			want: "a very long filename that exceeds the maximum allowed length of two hundred characters for obsidian vault file names and should be truncat",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeFilename(tt.in)
			if got != tt.want {
				t.Errorf("sanitizeFilename() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetContentType(t *testing.T) {
	tests := []struct {
		filename string
		want     string
	}{
		{"test.md", "text/markdown"},
		{"test.markdown", "text/markdown"},
		{"test.txt", "text/plain"},
		{"test.html", "text/html"},
		{"test.json", "application/json"},
		{"test.png", "image/png"},
		{"test.jpg", "image/jpeg"},
		{"test.pdf", "application/pdf"},
		{"test.unknown", "application/octet-stream"},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			got := getContentType(tt.filename)
			if got != tt.want {
				t.Errorf("getContentType(%q) = %q, want %q", tt.filename, got, tt.want)
			}
		})
	}
}

func TestIsPathWithinVault(t *testing.T) {
	tests := []struct {
		path      string
		vaultPath string
		want      bool
	}{
		{
			path:      "/vault/documents/test.md",
			vaultPath: "/vault",
			want:      true,
		},
		{
			path:      "/vault/documents/sub/test.md",
			vaultPath: "/vault",
			want:      true,
		},
		{
			path:      "/other/documents/test.md",
			vaultPath: "/vault",
			want:      false,
		},
		{
			path:      "/vault/../etc/passwd",
			vaultPath: "/vault",
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := isPathWithinVault(tt.path, tt.vaultPath)
			if got != tt.want {
				t.Errorf("isPathWithinVault(%q, %q) = %v, want %v", tt.path, tt.vaultPath, got, tt.want)
			}
		})
	}
}

func TestBuildFrontmatter(t *testing.T) {
	meta := map[string]string{
		"title":   "Test Document",
		"author":  "user123",
		"created": time.Now().Format(time.RFC3339),
	}

	got := buildFrontmatter(meta)
	
	// Check that it starts and ends correctly
	if !startsWith(got, "---\n") {
		t.Error("buildFrontmatter should start with '---'")
	}
	if !endsWith(got, "---") {
		t.Error("buildFrontmatter should end with '---'")
	}
	
	// Check that it contains the metadata
	if !contains(got, "title: Test Document") {
		t.Error("buildFrontmatter should contain title")
	}
	if !contains(got, "author: user123") {
		t.Error("buildFrontmatter should contain author")
	}
}

// Helper functions
func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

func endsWith(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsAt(s, substr, 0))
}

func containsAt(s, substr string, start int) bool {
	if start > len(s)-len(substr) {
		return false
	}
	if s[start:start+len(substr)] == substr {
		return true
	}
	return containsAt(s, substr, start+1)
}
