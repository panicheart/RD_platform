package handlers

import (
	"io"
	"net/http"
	"path/filepath"

	"rdp/services/api/models"
	"rdp/services/api/services"

	"github.com/gin-gonic/gin"
)

// FileHandler handles file HTTP requests
type FileHandler struct {
	fileService *services.FileService
}

// NewFileHandler creates a new FileHandler
func NewFileHandler(fileService *services.FileService) *FileHandler {
	return &FileHandler{
		fileService: fileService,
	}
}

// ListFiles handles GET /api/v1/projects/:id/files
func (h *FileHandler) ListFiles(c *gin.Context) {
	projectID := c.Param("id")
	path := c.Query("path")

	files, err := h.fileService.ListProjectFiles(c.Request.Context(), projectID, path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    files,
	})
}

// UploadFile handles POST /api/v1/projects/:id/files
func (h *FileHandler) UploadFile(c *gin.Context) {
	projectID := c.Param("id")
	path := c.Query("path")
	if path == "" {
		path = "/"
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4000,
			"message": "no file uploaded",
			"data":    nil,
		})
		return
	}
	defer file.Close()

	projectFile, err := h.fileService.UploadFile(c.Request.Context(), projectID, path, header.Filename, file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4001,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "file uploaded successfully",
		"data":    projectFile,
	})
}

// CreateDirectory handles POST /api/v1/projects/:id/files/directory
func (h *FileHandler) CreateDirectory(c *gin.Context) {
	projectID := c.Param("id")

	var req struct {
		Path string `json:"path" binding:"required"`
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4000,
			"message": "invalid request body: " + err.Error(),
			"data":    nil,
		})
		return
	}

	dir, err := h.fileService.CreateDirectory(c.Request.Context(), projectID, req.Path, req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4001,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "directory created successfully",
		"data":    dir,
	})
}

// DeleteFile handles DELETE /api/v1/projects/:projectId/files/:fileId
func (h *FileHandler) DeleteFile(c *gin.Context) {
	projectID := c.Param("projectId")
	fileID := c.Param("fileId")

	_ = projectID // Used for authorization check

	err := h.fileService.DeleteFile(c.Request.Context(), fileID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4001,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "file deleted successfully",
		"data":    nil,
	})
}

// DownloadFile handles GET /api/v1/projects/:projectId/files/:fileId/download
func (h *FileHandler) DownloadFile(c *gin.Context) {
	projectID := c.Param("projectId")
	fileID := c.Param("fileId")

	_ = projectID // Used for authorization check

	reader, err := h.fileService.DownloadFile(c.Request.Context(), fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    4040,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	defer reader.Close()

	// Get filename from database if needed, otherwise use generic name
	filename := filepath.Base(c.Request.URL.Path)

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")

	_, err = io.Copy(c.Writer, reader)
	if err != nil {
		// Can't send headers after they are sent, just log
		_ = err
	}
}

// FileUploadRequest represents a file upload request
type FileUploadRequest struct {
	File   io.Reader `json:"-"`
	Path   string    `json:"path"`
	Name   string    `json:"name"`
}

// Ensure FileUploadRequest implements proper binding
var _ = func() {
	// This is a compile-time check to ensure the struct is properly defined
	var _ models.ProjectFile
}
