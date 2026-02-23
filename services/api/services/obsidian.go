package services

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"gorm.io/gorm"
	"rdp-platform/rdp-api/models"
)

// ObsidianService provides WebDAV and sync functionality for Obsidian integration
type ObsidianService struct {
	db *gorm.DB
}

// NewObsidianService creates a new Obsidian service
func NewObsidianService(db *gorm.DB) *ObsidianService {
	return &ObsidianService{db: db}
}

// GetVaults returns all configured vaults for a user
func (s *ObsidianService) GetVaults(userID string) ([]models.ObsidianMapping, error) {
	var mappings []models.ObsidianMapping
	if err := s.db.Find(&mappings).Error; err != nil {
		return nil, err
	}
	return mappings, nil
}

// GetVault returns a single vault by ID
func (s *ObsidianService) GetVault(vaultID string) (*models.ObsidianMapping, error) {
	var mapping models.ObsidianMapping
	if err := s.db.First(&mapping, "id = ?", vaultID).Error; err != nil {
		return nil, err
	}
	return &mapping, nil
}

// CreateVault creates a new vault mapping
func (s *ObsidianService) CreateVault(vaultPath, localPath, categoryID string, autoSync bool) (*models.ObsidianMapping, error) {
	mapping := &models.ObsidianMapping{
		VaultPath:  vaultPath,
		LocalPath:  localPath,
		CategoryID: categoryID,
		AutoSync:   autoSync,
	}

	if err := s.db.Create(mapping).Error; err != nil {
		return nil, err
	}

	return mapping, nil
}

// UpdateVault updates a vault mapping
func (s *ObsidianService) UpdateVault(vaultID string, vaultPath, localPath, categoryID string, autoSync bool) (*models.ObsidianMapping, error) {
	var mapping models.ObsidianMapping
	if err := s.db.First(&mapping, "id = ?", vaultID).Error; err != nil {
		return nil, err
	}

	mapping.VaultPath = vaultPath
	mapping.LocalPath = localPath
	mapping.CategoryID = categoryID
	mapping.AutoSync = autoSync

	if err := s.db.Save(&mapping).Error; err != nil {
		return nil, err
	}

	return &mapping, nil
}

// DeleteVault deletes a vault mapping
func (s *ObsidianService) DeleteVault(vaultID string) error {
	return s.db.Delete(&models.ObsidianMapping{}, "id = ?", vaultID).Error
}

// WebDAVProps represents WebDAV resource properties
type WebDAVProps struct {
	Path          string    `xml:"href"`
	DisplayName   string    `xml:"displayname"`
	ContentType   string    `xml:"getcontenttype"`
	ContentLength int64     `xml:"getcontentlength"`
	LastModified  time.Time `xml:"getlastmodified"`
	IsCollection  bool      `xml:"iscollection"`
}

// PropFind performs WebDAV PROPFIND operation
func (s *ObsidianService) PropFind(vaultID, path string, depth int) ([]WebDAVProps, error) {
	mapping, err := s.GetVault(vaultID)
	if err != nil {
		return nil, fmt.Errorf("vault not found: %w", err)
	}

	fullPath := filepath.Join(mapping.LocalPath, path)

	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, fmt.Errorf("path not found: %w", err)
	}

	props := []WebDAVProps{
		{
			Path:         path,
			DisplayName:  filepath.Base(path),
			LastModified: info.ModTime(),
			IsCollection: info.IsDir(),
		},
	}

	if depth > 0 && info.IsDir() {
		entries, err := os.ReadDir(fullPath)
		if err != nil {
			return nil, err
		}

		for _, entry := range entries {
			entryPath := filepath.Join(path, entry.Name())
			entryInfo, err := entry.Info()
			if err != nil {
				continue
			}

			entryProps := WebDAVProps{
				Path:         entryPath,
				DisplayName:  entry.Name(),
				LastModified: entryInfo.ModTime(),
				IsCollection: entry.IsDir(),
			}

			if !entry.IsDir() {
				entryProps.ContentLength = entryInfo.Size()
				entryProps.ContentType = getContentType(entry.Name())
			}

			props = append(props, entryProps)
		}
	}

	return props, nil
}

// GetFile retrieves a file from the vault
func (s *ObsidianService) GetFile(vaultID, path string) (io.ReadCloser, int64, string, error) {
	mapping, err := s.GetVault(vaultID)
	if err != nil {
		return nil, 0, "", fmt.Errorf("vault not found: %w", err)
	}

	fullPath := filepath.Join(mapping.LocalPath, path)

	if !isPathWithinVault(fullPath, mapping.LocalPath) {
		return nil, 0, "", fmt.Errorf("access denied: path outside vault")
	}

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, 0, "", err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, 0, "", err
	}

	if info.IsDir() {
		file.Close()
		return nil, 0, "", fmt.Errorf("path is a directory")
	}

	contentType := getContentType(path)
	return file, info.Size(), contentType, nil
}

// PutFile stores a file in the vault
func (s *ObsidianService) PutFile(vaultID, path string, content io.Reader, size int64) error {
	mapping, err := s.GetVault(vaultID)
	if err != nil {
		return fmt.Errorf("vault not found: %w", err)
	}

	fullPath := filepath.Join(mapping.LocalPath, path)

	if !isPathWithinVault(fullPath, mapping.LocalPath) {
		return fmt.Errorf("access denied: path outside vault")
	}

	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	tempPath := fullPath + ".tmp"
	file, err := os.Create(tempPath)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, content)
	file.Close()

	if err != nil {
		os.Remove(tempPath)
		return err
	}

	if err := os.Rename(tempPath, fullPath); err != nil {
		os.Remove(tempPath)
		return err
	}

	return nil
}

// DeleteFile deletes a file or directory from the vault
func (s *ObsidianService) DeleteFile(vaultID, path string) error {
	mapping, err := s.GetVault(vaultID)
	if err != nil {
		return fmt.Errorf("vault not found: %w", err)
	}

	fullPath := filepath.Join(mapping.LocalPath, path)

	if !isPathWithinVault(fullPath, mapping.LocalPath) {
		return fmt.Errorf("access denied: path outside vault")
	}

	return os.RemoveAll(fullPath)
}

// MkCol creates a directory in the vault
func (s *ObsidianService) MkCol(vaultID, path string) error {
	mapping, err := s.GetVault(vaultID)
	if err != nil {
		return fmt.Errorf("vault not found: %w", err)
	}

	fullPath := filepath.Join(mapping.LocalPath, path)

	if !isPathWithinVault(fullPath, mapping.LocalPath) {
		return fmt.Errorf("access denied: path outside vault")
	}

	return os.MkdirAll(fullPath, 0755)
}

// Move moves a file or directory within the vault
func (s *ObsidianService) Move(vaultID, sourcePath, destPath string) error {
	mapping, err := s.GetVault(vaultID)
	if err != nil {
		return fmt.Errorf("vault not found: %w", err)
	}

	fullSource := filepath.Join(mapping.LocalPath, sourcePath)
	fullDest := filepath.Join(mapping.LocalPath, destPath)

	if !isPathWithinVault(fullSource, mapping.LocalPath) || !isPathWithinVault(fullDest, mapping.LocalPath) {
		return fmt.Errorf("access denied: path outside vault")
	}

	destDir := filepath.Dir(fullDest)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	return os.Rename(fullSource, fullDest)
}

// SyncResult represents the result of a sync operation
type SyncResult struct {
	Uploaded   int      `json:"uploaded"`
	Downloaded int      `json:"downloaded"`
	Conflicts  int      `json:"conflicts"`
	Errors     []string `json:"errors,omitempty"`
}

// SyncVault performs bidirectional sync between platform and Obsidian vault
func (s *ObsidianService) SyncVault(vaultID string) (*SyncResult, error) {
	mapping, err := s.GetVault(vaultID)
	if err != nil {
		return nil, fmt.Errorf("vault not found: %w", err)
	}

	result := &SyncResult{}

	var knowledgeList []models.Knowledge
	if err := s.db.Where("category_id = ? AND source = ?", mapping.CategoryID, "obsidian").Find(&knowledgeList).Error; err != nil {
		return nil, err
	}

	existingMap := make(map[string]*models.Knowledge)
	for i := range knowledgeList {
		existingMap[knowledgeList[i].SourceID] = &knowledgeList[i]
	}

	err = filepath.Walk(mapping.LocalPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("walk error at %s: %v", path, err))
			return nil
		}

		if info.IsDir() || !strings.HasSuffix(strings.ToLower(path), ".md") {
			return nil
		}

		relPath, err := filepath.Rel(mapping.LocalPath, path)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("path error: %v", err))
			return nil
		}

		relPath = filepath.ToSlash(relPath)

		if existing, found := existingMap[relPath]; found {
			if info.ModTime().After(existing.UpdatedAt) {
				if err := s.importMarkdownFile(path, relPath, existing, mapping.CategoryID); err != nil {
					result.Errors = append(result.Errors, fmt.Sprintf("import error for %s: %v", relPath, err))
				} else {
					result.Uploaded++
				}
			}
			delete(existingMap, relPath)
		} else {
			if err := s.importMarkdownFile(path, relPath, nil, mapping.CategoryID); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("import error for %s: %v", relPath, err))
			} else {
				result.Uploaded++
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	for _, orphaned := range existingMap {
		if err := s.exportToVault(orphaned, mapping.LocalPath); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("export error for %s: %v", orphaned.SourceID, err))
		} else {
			result.Downloaded++
		}
	}

	now := time.Now()
	mapping.LastSyncAt = &now
	s.db.Save(mapping)

	return result, nil
}

// ImportMarkdown imports a markdown file into the knowledge base
func (s *ObsidianService) ImportMarkdown(vaultID, filePath, content, authorID string) (*models.Knowledge, error) {
	mapping, err := s.GetVault(vaultID)
	if err != nil {
		return nil, fmt.Errorf("vault not found: %w", err)
	}

	metadata, markdownContent := parseFrontmatter(content)
	tags := extractTags(markdownContent)

	var existing models.Knowledge
	err = s.db.Where("source_id = ? AND source = ?", filePath, "obsidian").First(&existing).Error

	knowledge := &models.Knowledge{
		Title:      metadata["title"],
		Content:    markdownContent,
		CategoryID: mapping.CategoryID,
		AuthorID:   authorID,
		Source:     "obsidian",
		SourceID:   filePath,
		Status:     "draft",
	}

	if tagStr, ok := metadata["tags"]; ok {
		frontmatterTags := parseTags(tagStr)
		tags = mergeTags(tags, frontmatterTags)
	}

	if err == nil {
		knowledge.ID = existing.ID
		if err := s.db.Model(&existing).Updates(knowledge).Error; err != nil {
			return nil, err
		}
		s.syncTags(&existing, tags)
	} else {
		if err := s.db.Create(knowledge).Error; err != nil {
			return nil, err
		}
		s.syncTags(knowledge, tags)
	}

	return knowledge, nil
}

// ExportToVault exports a knowledge entry to Obsidian vault
func (s *ObsidianService) ExportToVault(knowledgeID, vaultID string) error {
	var knowledge models.Knowledge
	if err := s.db.First(&knowledge, "id = ?", knowledgeID).Error; err != nil {
		return err
	}

	mapping, err := s.GetVault(vaultID)
	if err != nil {
		return fmt.Errorf("vault not found: %w", err)
	}

	return s.exportToVault(&knowledge, mapping.LocalPath)
}

// importMarkdownFile imports a single markdown file from vault
func (s *ObsidianService) importMarkdownFile(fullPath, relPath string, existing *models.Knowledge, categoryID string) error {
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return err
	}

	metadata, markdownContent := parseFrontmatter(string(content))
	tags := extractTags(markdownContent)

	knowledge := &models.Knowledge{
		Title:      metadata["title"],
		Content:    markdownContent,
		CategoryID: categoryID,
		Source:     "obsidian",
		SourceID:   relPath,
		Status:     "draft",
	}

	if metadata["author"] != "" {
		knowledge.AuthorID = metadata["author"]
	}

	if tagStr, ok := metadata["tags"]; ok {
		frontmatterTags := parseTags(tagStr)
		tags = mergeTags(tags, frontmatterTags)
	}

	if existing != nil {
		knowledge.ID = existing.ID
		if err := s.db.Model(existing).Updates(knowledge).Error; err != nil {
			return err
		}
		s.syncTags(existing, tags)
	} else {
		if err := s.db.Create(knowledge).Error; err != nil {
			return err
		}
		s.syncTags(knowledge, tags)
	}

	return nil
}

// exportToVault exports a knowledge entry to the vault
func (s *ObsidianService) exportToVault(knowledge *models.Knowledge, vaultPath string) error {
	filePath := knowledge.SourceID
	if filePath == "" {
		safeTitle := sanitizeFilename(knowledge.Title)
		filePath = safeTitle + ".md"
	}

	fullPath := filepath.Join(vaultPath, filePath)

	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	var tags []models.Tag
	s.db.Model(knowledge).Association("Tags").Find(&tags)
	tagNames := make([]string, len(tags))
	for i, tag := range tags {
		tagNames[i] = tag.Name
	}

	frontmatter := buildFrontmatter(map[string]string{
		"title":   knowledge.Title,
		"id":      knowledge.ID,
		"author":  knowledge.AuthorID,
		"tags":    strings.Join(tagNames, ", "),
		"created": knowledge.CreatedAt.Format(time.RFC3339),
		"updated": knowledge.UpdatedAt.Format(time.RFC3339),
	})

	content := frontmatter + "\n" + knowledge.Content
	return os.WriteFile(fullPath, []byte(content), 0644)
}

// syncTags synchronizes tags for a knowledge entry
func (s *ObsidianService) syncTags(knowledge *models.Knowledge, tagNames []string) error {
	if err := s.db.Model(knowledge).Association("Tags").Clear(); err != nil {
		return err
	}

	for _, tagName := range tagNames {
		var tag models.Tag
		err := s.db.FirstOrCreate(&tag, models.Tag{Name: tagName}).Error
		if err != nil {
			continue
		}
		s.db.Model(knowledge).Association("Tags").Append(&tag)
	}

	return nil
}

// parseFrontmatter parses YAML frontmatter from markdown content
func parseFrontmatter(content string) (map[string]string, string) {
	metadata := make(map[string]string)

	if !strings.HasPrefix(content, "---\n") {
		return metadata, content
	}

	parts := strings.SplitN(content, "---\n", 3)
	if len(parts) < 3 {
		return metadata, content
	}

	lines := strings.Split(parts[1], "\n")
	for _, line := range lines {
		if strings.Contains(line, ":") {
			colonIdx := strings.Index(line, ":")
			key := strings.TrimSpace(line[:colonIdx])
			value := strings.TrimSpace(line[colonIdx+1:])
			value = strings.Trim(value, `"'`)
			metadata[key] = value
		}
	}

	return metadata, strings.TrimSpace(parts[2])
}

// buildFrontmatter creates YAML frontmatter from metadata
func buildFrontmatter(metadata map[string]string) string {
	var sb strings.Builder
	sb.WriteString("---\n")
	for key, value := range metadata {
		if strings.Contains(value, ":") || strings.Contains(value, "\n") {
			sb.WriteString(fmt.Sprintf("%s: \"%s\"\n", key, value))
		} else {
			sb.WriteString(fmt.Sprintf("%s: %s\n", key, value))
		}
	}
	sb.WriteString("---")
	return sb.String()
}

// extractTags extracts #tags from markdown content
func extractTags(content string) []string {
	tagRegex := regexp.MustCompile(`#([\w\-\u4e00-\u9fa5]+)`)
	matches := tagRegex.FindAllStringSubmatch(content, -1)

	tags := make([]string, 0, len(matches))
	seen := make(map[string]bool)

	for _, match := range matches {
		if len(match) > 1 && !seen[match[1]] {
			tags = append(tags, match[1])
			seen[match[1]] = true
		}
	}

	return tags
}

// parseTags parses comma-separated tag string
func parseTags(tagStr string) []string {
	parts := strings.Split(tagStr, ",")
	tags := make([]string, 0, len(parts))

	for _, part := range parts {
		tag := strings.TrimSpace(part)
		if tag != "" {
			tags = append(tags, tag)
		}
	}

	return tags
}

// mergeTags merges two tag slices without duplicates
func mergeTags(a, b []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(a)+len(b))

	for _, tag := range a {
		if !seen[tag] {
			result = append(result, tag)
			seen[tag] = true
		}
	}

	for _, tag := range b {
		if !seen[tag] {
			result = append(result, tag)
			seen[tag] = true
		}
	}

	return result
}

// sanitizeFilename sanitizes a string for use as filename
func sanitizeFilename(name string) string {
	invalid := regexp.MustCompile(`[<>:"/\\|?*]`)
	name = invalid.ReplaceAllString(name, "_")
	name = strings.TrimSpace(name)
	name = strings.Trim(name, ".")
	if len(name) > 200 {
		name = name[:200]
	}
	return name
}

// isPathWithinVault checks if a path is within the vault directory
func isPathWithinVault(path, vaultPath string) bool {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	absVault, err := filepath.Abs(vaultPath)
	if err != nil {
		return false
	}

	return strings.HasPrefix(absPath, absVault)
}

// getContentType returns MIME type for a file
func getContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".md", ".markdown":
		return "text/markdown"
	case ".txt":
		return "text/plain"
	case ".html", ".htm":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".json":
		return "application/json"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".pdf":
		return "application/pdf"
	default:
		return "application/octet-stream"
	}
}

// HandleWebDAVRequest handles incoming WebDAV requests
func (s *ObsidianService) HandleWebDAVRequest(vaultID string, r *http.Request) (int, interface{}, error) {
	switch r.Method {
	case "OPTIONS":
		return http.StatusOK, nil, nil

	case "PROPFIND":
		depth := 0
		if d := r.Header.Get("Depth"); d == "1" {
			depth = 1
		} else if d == "infinity" {
			depth = 1
		}

		props, err := s.PropFind(vaultID, r.URL.Path, depth)
		if err != nil {
			return http.StatusNotFound, nil, err
		}
		return http.StatusMultiStatus, props, nil

	case "GET", "HEAD":
		reader, size, contentType, err := s.GetFile(vaultID, r.URL.Path)
		if err != nil {
			return http.StatusNotFound, nil, err
		}
		defer reader.Close()

		if r.Method == "HEAD" {
			return http.StatusOK, map[string]interface{}{
				"content_type":   contentType,
				"content_length": size,
			}, nil
		}

		content, _ := io.ReadAll(reader)
		return http.StatusOK, map[string]interface{}{
			"content":      string(content),
			"content_type": contentType,
			"size":         size,
		}, nil

	case "PUT":
		content := r.Body
		defer content.Close()

		if err := s.PutFile(vaultID, r.URL.Path, content, r.ContentLength); err != nil {
			return http.StatusInternalServerError, nil, err
		}
		return http.StatusCreated, nil, nil

	case "DELETE":
		if err := s.DeleteFile(vaultID, r.URL.Path); err != nil {
			return http.StatusNotFound, nil, err
		}
		return http.StatusNoContent, nil, nil

	case "MKCOL":
		if err := s.MkCol(vaultID, r.URL.Path); err != nil {
			return http.StatusConflict, nil, err
		}
		return http.StatusCreated, nil, nil

	case "MOVE":
		dest := r.Header.Get("Destination")
		if dest == "" {
			return http.StatusBadRequest, nil, fmt.Errorf("missing Destination header")
		}

		destPath := dest
		if idx := strings.Index(dest, "/webdav/"); idx != -1 {
			parts := strings.SplitN(dest[idx+8:], "/", 2)
			if len(parts) > 1 {
				destPath = "/" + parts[1]
			}
		}

		if err := s.Move(vaultID, r.URL.Path, destPath); err != nil {
			return http.StatusConflict, nil, err
		}
		return http.StatusCreated, nil, nil

	default:
		return http.StatusMethodNotAllowed, nil, fmt.Errorf("method not allowed: %s", r.Method)
	}
}
