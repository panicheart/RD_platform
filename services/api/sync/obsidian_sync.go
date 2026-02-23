package sync

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"
	"rdp-platform/rdp-api/models"
)

// ObsidianSync provides bidirectional synchronization between Obsidian vault and platform
type ObsidianSync struct {
	db *gorm.DB
}

// NewObsidianSync creates a new Obsidian sync instance
func NewObsidianSync(db *gorm.DB) *ObsidianSync {
	return &ObsidianSync{db: db}
}

// SyncConfig represents configuration for a sync operation
type SyncConfig struct {
	VaultID       string
	Direction     string // "bidirectional", "upload", "download"
	ConflictPolicy string // "platform_wins", "vault_wins", "newer_wins"
	IncludeFiles  []string
	ExcludeFiles  []string
}

// SyncStats represents statistics for a sync operation
type SyncStats struct {
	FilesProcessed int
	FilesUploaded  int
	FilesDownloaded int
	FilesConflicted int
	FilesSkipped   int
	Errors         []string
}

// PerformSync executes the sync operation
func (s *ObsidianSync) PerformSync(config SyncConfig) (*SyncStats, error) {
	// Get vault mapping
	var mapping models.ObsidianMapping
	if err := s.db.First(&mapping, "id = ?", config.VaultID).Error; err != nil {
		return nil, fmt.Errorf("vault not found: %w", err)
	}

	stats := &SyncStats{}

	// Get all knowledge entries for this category
	var knowledgeList []models.Knowledge
	if err := s.db.Where("category_id = ? AND source = ?", mapping.CategoryID, "obsidian").Find(&knowledgeList).Error; err != nil {
		return nil, err
	}

	// Build map of platform entries
	platformFiles := make(map[string]*models.Knowledge)
	for i := range knowledgeList {
		platformFiles[knowledgeList[i].SourceID] = &knowledgeList[i]
	}

	// Scan vault directory
	vaultFiles := make(map[string]os.FileInfo)
	err := filepath.Walk(mapping.LocalPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			stats.Errors = append(stats.Errors, fmt.Sprintf("walk error: %v", err))
			return nil
		}

		if info.IsDir() || !s.shouldSyncFile(path, config) {
			return nil
		}

		relPath, _ := filepath.Rel(mapping.LocalPath, path)
		relPath = filepath.ToSlash(relPath)
		vaultFiles[relPath] = info

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Process files based on direction
	switch config.Direction {
	case "upload":
		s.syncUpload(platformFiles, vaultFiles, mapping.LocalPath, config, stats)
	case "download":
		s.syncDownload(platformFiles, vaultFiles, mapping.LocalPath, config, stats)
	default: // bidirectional
		s.syncBidirectional(platformFiles, vaultFiles, mapping.LocalPath, config, stats)
	}

	// Update last sync time
	now := time.Now()
	mapping.LastSyncAt = &now
	s.db.Save(&mapping)

	return stats, nil
}

// syncUpload syncs from vault to platform (upload only)
func (s *ObsidianSync) syncUpload(platformFiles map[string]*models.Knowledge, vaultFiles map[string]os.FileInfo, vaultPath string, config SyncConfig, stats *SyncStats) {
	for relPath := range vaultFiles {
		stats.FilesProcessed++

		if _, exists := platformFiles[relPath]; !exists {
			// New file in vault - upload to platform
			if err := s.importFile(vaultPath, relPath, config); err != nil {
				stats.Errors = append(stats.Errors, fmt.Sprintf("import %s: %v", relPath, err))
			} else {
				stats.FilesUploaded++
			}
		}
	}
}

// syncDownload syncs from platform to vault (download only)
func (s *ObsidianSync) syncDownload(platformFiles map[string]*models.Knowledge, vaultFiles map[string]os.FileInfo, vaultPath string, config SyncConfig, stats *SyncStats) {
	for relPath, knowledge := range platformFiles {
		stats.FilesProcessed++

		if _, exists := vaultFiles[relPath]; !exists {
			// New file in platform - download to vault
			if err := s.exportFile(knowledge, vaultPath); err != nil {
				stats.Errors = append(stats.Errors, fmt.Sprintf("export %s: %v", relPath, err))
			} else {
				stats.FilesDownloaded++
			}
		}
	}
}

// syncBidirectional syncs in both directions with conflict resolution
func (s *ObsidianSync) syncBidirectional(platformFiles map[string]*models.Knowledge, vaultFiles map[string]os.FileInfo, vaultPath string, config SyncConfig, stats *SyncStats) {
	allFiles := make(map[string]bool)
	for path := range platformFiles {
		allFiles[path] = true
	}
	for path := range vaultFiles {
		allFiles[path] = true
	}

	for relPath := range allFiles {
		stats.FilesProcessed++

		platformEntry, inPlatform := platformFiles[relPath]
		vaultInfo, inVault := vaultFiles[relPath]

		if inPlatform && !inVault {
			// Only in platform - download
			if err := s.exportFile(platformEntry, vaultPath); err != nil {
				stats.Errors = append(stats.Errors, fmt.Sprintf("export %s: %v", relPath, err))
			} else {
				stats.FilesDownloaded++
			}
		} else if inVault && !inPlatform {
			// Only in vault - upload
			if err := s.importFile(vaultPath, relPath, config); err != nil {
				stats.Errors = append(stats.Errors, fmt.Sprintf("import %s: %v", relPath, err))
			} else {
				stats.FilesUploaded++
			}
		} else {
			// Exists in both - check for conflict
			if vaultInfo.ModTime().After(platformEntry.UpdatedAt) {
				// Vault is newer
				switch config.ConflictPolicy {
				case "vault_wins", "newer_wins":
					if err := s.importFile(vaultPath, relPath, config); err != nil {
						stats.Errors = append(stats.Errors, fmt.Sprintf("import %s: %v", relPath, err))
					} else {
						stats.FilesUploaded++
					}
				default: // platform_wins
					stats.FilesSkipped++
				}
			} else if platformEntry.UpdatedAt.After(vaultInfo.ModTime()) {
				// Platform is newer
				switch config.ConflictPolicy {
				case "platform_wins", "newer_wins":
					if err := s.exportFile(platformEntry, vaultPath); err != nil {
						stats.Errors = append(stats.Errors, fmt.Sprintf("export %s: %v", relPath, err))
					} else {
						stats.FilesDownloaded++
					}
				default: // vault_wins
					stats.FilesSkipped++
				}
			} else {
				stats.FilesSkipped++
			}
			stats.FilesConflicted++
		}
	}
}

// importFile imports a file from vault to platform
func (s *ObsidianSync) importFile(vaultPath, relPath string, config SyncConfig) error {
	fullPath := filepath.Join(vaultPath, relPath)
	
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return err
	}

	// Get mapping for category
	var mapping models.ObsidianMapping
	if err := s.db.Where("local_path = ?", vaultPath).First(&mapping).Error; err != nil {
		return err
	}

	// Parse YAML frontmatter
	metadata, markdownContent := parseFrontmatter(string(content))

	// Check if file already exists
	var existing models.Knowledge
	err = s.db.Where("source_id = ? AND source = ?", relPath, "obsidian").First(&existing).Error

	knowledge := &models.Knowledge{
		Title:      metadata["title"],
		Content:    markdownContent,
		CategoryID: mapping.CategoryID,
		Source:     "obsidian",
		SourceID:   relPath,
		Status:     "draft",
	}

	if metadata["author"] != "" {
		knowledge.AuthorID = metadata["author"]
	} else {
		knowledge.AuthorID = "system"
	}

	// Handle tags
	tags := extractTags(markdownContent)
	if tagStr, ok := metadata["tags"]; ok {
		tags = mergeTags(tags, parseTags(tagStr))
	}

	if err == nil {
		// Update existing
		knowledge.ID = existing.ID
		if err := s.db.Model(&existing).Updates(knowledge).Error; err != nil {
			return err
		}
		s.syncTags(&existing, tags)
	} else {
		// Create new
		if err := s.db.Create(knowledge).Error; err != nil {
			return err
		}
		s.syncTags(knowledge, tags)
	}

	return nil
}

// exportFile exports a file from platform to vault
func (s *ObsidianSync) exportFile(knowledge *models.Knowledge, vaultPath string) error {
	filePath := knowledge.SourceID
	if filePath == "" {
		filePath = sanitizeFilename(knowledge.Title) + ".md"
	}

	fullPath := filepath.Join(vaultPath, filePath)
	
	// Ensure directory exists
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Get tags
	var tags []models.Tag
	s.db.Model(knowledge).Association("Tags").Find(&tags)
	tagNames := make([]string, len(tags))
	for i, tag := range tags {
		tagNames[i] = tag.Name
	}

	// Build YAML frontmatter
	frontmatter := buildFrontmatter(map[string]string{
		"title":   knowledge.Title,
		"id":      knowledge.ID,
		"author":  knowledge.AuthorID,
		"tags":    strings.Join(tagNames, ", "),
		"created": knowledge.CreatedAt.Format(time.RFC3339),
		"updated": knowledge.UpdatedAt.Format(time.RFC3339),
	})

	// Write file
	content := frontmatter + "\n" + knowledge.Content
	return os.WriteFile(fullPath, []byte(content), 0644)
}

// syncTags synchronizes tags for a knowledge entry
func (s *ObsidianSync) syncTags(knowledge *models.Knowledge, tagNames []string) error {
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

// shouldSyncFile checks if a file should be synced based on filters
func (s *ObsidianSync) shouldSyncFile(path string, config SyncConfig) bool {
	// Check if it's a markdown file
	if !strings.HasSuffix(strings.ToLower(path), ".md") {
		return false
	}

	// Check exclude patterns
	for _, pattern := range config.ExcludeFiles {
		if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
			return false
		}
	}

	// Check include patterns (if specified)
	if len(config.IncludeFiles) > 0 {
		included := false
		for _, pattern := range config.IncludeFiles {
			if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
				included = true
				break
			}
		}
		return included
	}

	return true
}

// parseFrontmatter parses YAML frontmatter from markdown
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

// buildFrontmatter creates YAML frontmatter
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

// extractTags extracts #tags from markdown
func extractTags(content string) []string {
	var tags []string
	seen := make(map[string]bool)

	words := strings.Fields(content)
	for _, word := range words {
		if strings.HasPrefix(word, "#") {
			tag := strings.TrimPrefix(word, "#")
			tag = strings.TrimFunc(tag, func(r rune) bool {
				return r == '#' || r == ' ' || r == '\t' || r == '\n'
			})
			if tag != "" && !seen[tag] {
				tags = append(tags, tag)
				seen[tag] = true
			}
		}
	}

	return tags
}

// parseTags parses comma-separated tags
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

// mergeTags merges tag slices without duplicates
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
	invalid := []string{
		"<", ">", ":", "\"", "/", "\\", "|", "?", "*",
	}

	for _, char := range invalid {
		name = strings.ReplaceAll(name, char, "_")
	}

	name = strings.TrimSpace(name)
	name = strings.Trim(name, ".")

	if len(name) > 200 {
		name = name[:200]
	}

	return name
}
