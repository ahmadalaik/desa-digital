package helpers

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/ahmadalaik/desa-digital/structs"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadConfig struct {
	File           *multipart.FileHeader
	AllowedTypes   []string
	MaxSize        int64
	DestinationDir string
}

type UploadResult struct {
	FileName string
	FilePath string
	Error    error
	Response *structs.ErrorResponse
}

func SlugifyFileName(filename string) string {
	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)

	slugBase := Slugify(base)

	return slugBase + ext
}

func UploadFile(c *gin.Context, config UploadConfig) UploadResult {
	if config.File == nil {
		return UploadResult{
			Response: &structs.ErrorResponse{
				Success: false,
				Message: "File is required",
				Errors:  map[string]string{"file": "No file was uploaded"},
			},
		}
	}

	if config.File.Size > config.MaxSize {
		return UploadResult{
			Response: &structs.ErrorResponse{
				Success: false,
				Message: "File size to large",
				Errors:  map[string]string{"file": fmt.Sprintf("Maximum file is: %dMB", config.MaxSize/(1<<20))},
			},
		}
	}

	ext := strings.ToLower(filepath.Ext(config.File.Filename))
	allowed := false

	// for _, t := range config.AllowedTypes {
	// 	if ext == t {
	// 		allowed = true
	// 		break
	// 	}
	// }

	allowed = slices.Contains(config.AllowedTypes, ext)

	if !allowed {
		return UploadResult{
			Response: &structs.ErrorResponse{
				Success: false,
				Message: "Invalid file type",
				Errors:  map[string]string{"file": fmt.Sprintf("Allowed file types: %v", config.AllowedTypes)},
			},
		}
	}

	uuidName := uuid.New().String()

	fileName := uuidName + ext
	filePath := filepath.Join(config.DestinationDir, fileName)

	if err := os.MkdirAll(config.DestinationDir, 0755); err != nil {
		return UploadResult{
			Response: &structs.ErrorResponse{
				Success: false,
				Message: "Failed to create upload directory",
				Errors:  map[string]string{"system": err.Error()},
			},
		}
	}

	if err := c.SaveUploadedFile(config.File, filePath); err != nil {
		return UploadResult{
			Response: &structs.ErrorResponse{
				Success: false,
				Message: "Failed to save file",
				Errors:  map[string]string{"file": err.Error()},
			},
		}
	}

	return UploadResult{
		FileName: fileName,
		FilePath: filePath,
	}
}
