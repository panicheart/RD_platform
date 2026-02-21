package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"rdp/services/api/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FileService handles file management business logic
type FileService struct {
	db        *gorm.DB
	basePath  string
}

// NewFileService creates a new FileService
func NewFileService(db *gorm.DB, basePath string) *FileService {
	return &FileService{
		db:        db,
		basePath:  basePath,
	}
}

// ListProjectFiles returns files for a project
func (s *FileService) ListProjectFiles(ctx context.Context, projectID string, path string) ([]models.File, error) {
	var files []models.ProjectFile
	projectUID, err := uuid.Parse(projectID)
	if err != nil {
		return nil, errors.New("invalid project ID")
	}

	query := s.db.Model(&models.ProjectFile{}).Where("project_id = ?", projectUID)

	if path != "" {
		query = query.Where("path = ?", path)
	} else {
		// List root files if no path specified
		query = query.Where("path = ?", "/")
	}

	if err := query.Order("is_directory DESC, name ASC").Find(&files).Error; err != nil {
		return nil, err
	}

	// Convert to generic File type
	result := make([]models.File, len(files))
	for i, f := range files {
		result[i] = models.File{
			ID:          f.ID,
			Name:        f.Name,
			Path:        f.Path,
			Size:        f.Size,
			ContentType: f.ContentType,
			IsDirectory: f.IsDirectory,
			CreatedAt:   f.CreatedAt,
			UpdatedAt:   f.UpdatedAt,
		}
	}

	return result, nil
}

// UploadFile uploads a file to a project
func (s *FileService) UploadFile(ctx context.Context, projectID, path, filename string, reader io.Reader) (*models.ProjectFile, error) {
	projectUID, err := uuid.Parse(projectID)
	if err != nil {
		return nil, errors.New("invalid project ID")
	}

	// Generate file ID
	fileID := uuid.New()

	// Create directory path
	dirPath := filepath.Join(s.basePath, projectID, path)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Create file path
	filePath := filepath.Join(dirPath, filename)

	// Copy file content
	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	written, err := io.Copy(file, reader)
	if err != nil {
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	// Determine content type
	contentType := "application/octet-stream"

	// Save to database
	projectFile := models.ProjectFile{
		ID:          fileID,
		ProjectID:   projectUID,
		Name:        filename,
		Path:        path,
		Size:        written,
		ContentType: contentType,
		IsDirectory: false,
		StoragePath: filePath,
	}

	if err := s.db.Create(&projectFile).Error; err != nil {
		os.Remove(filePath)
		return nil, err
	}

	return &projectFile, nil
}

// CreateDirectory creates a directory in a project
func (s *FileService) CreateDirectory(ctx context.Context, projectID, path, name string) (*models.ProjectFile, error) {
	projectUID, err := uuid.Parse(projectID)
	if err != nil {
		return nil, errors.New("invalid project ID")
	}

	// Check if directory already exists
	var existing models.ProjectFile
	if err := s.db.First(&existing, "project_id = ? AND path = ? AND name = ? AND is_directory = ?", projectUID, path, name, true).Error; err == nil {
		return nil, errors.New("directory already exists")
	}

	dirID := uuid.New()

	// Create physical directory
	dirPath := filepath.Join(s.basePath, projectID, path, name)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Save to database
	projectDir := models.ProjectFile{
		ID:          dirID,
		ProjectID:   projectUID,
		Name:        name,
		Path:        path,
		IsDirectory: true,
		StoragePath: dirPath,
	}

	if err := s.db.Create(&projectDir).Error; err != nil {
		os.Remove(dirPath)
		return nil, err
	}

	return &projectDir, nil
}

// DeleteFile deletes a file or directory
func (s *FileService) DeleteFile(ctx context.Context, fileID string) error {
	uid, err := uuid.Parse(fileID)
	if err != nil {
		return errors.New("invalid file ID")
	}

	var file models.ProjectFile
	if err := s.db.First(&file, "id = ?", uid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("file not found")
		}
		return err
	}

	// Delete physical file/directory
	if file.IsDirectory {
		if err := os.RemoveAll(file.StoragePath); err != nil {
			return fmt.Errorf("failed to delete directory: %w", err)
		}
	} else {
		if err := os.Remove(file.StoragePath); err != nil {
			return fmt.Errorf("failed to delete file: %w", err)
		}
	}

	// Delete from database
	if err := s.db.Delete(&file).Error; err != nil {
		return err
	}

	return nil
}

// DownloadFile returns a reader for a file
func (s *FileService) DownloadFile(ctx context.Context, fileID string) (io.ReadCloser, error) {
	uid, err := uuid.Parse(fileID)
	if err != nil {
		return nil, errors.New("invalid file ID")
	}

	var file models.ProjectFile
	if err := s.db.First(&file, "id = ? AND is_directory = ?", uid, false).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("file not found")
		}
		return nil, err
	}

	reader, err := os.Open(file.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return reader, nil
}
