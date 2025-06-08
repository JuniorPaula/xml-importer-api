package handlers

import (
	"fmt"
	"importerapi/internal/models"
	"importerapi/internal/repositories"
	"importerapi/internal/worker"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ImportHandler struct {
	JobQueue         chan worker.ImportJob
	ImportStatusRepo repositories.ImportStatusRepository
}

func NewImportHandler(queue chan worker.ImportJob, db *gorm.DB) *ImportHandler {
	return &ImportHandler{
		JobQueue:         queue,
		ImportStatusRepo: repositories.NewImportStatusRepo(db),
	}
}

// ImportXMLDataHandler handles the import of data from a XML file
func (h *ImportHandler) ImportXMLDataHandler(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to get file from request: " + err.Error(),
		})
	}

	err = os.Mkdir("uploads", os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create uploads directory: " + err.Error(),
		})
	}

	filename := fmt.Sprintf("uploads/%d_%s", time.Now().Unix(), fileHeader.Filename)
	if err := c.SaveFile(fileHeader, filename); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to save file: " + err.Error(),
		})
	}

	// Generate a random import ID for tracking
	importID := uuid.New().String()
	err = h.ImportStatusRepo.Create(&models.ImportStatus{
		ImportID: importID,
		FileName: fileHeader.Filename,
		Status:   "queued",
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create import status: " + err.Error(),
		})
	}

	// Add the job to the queue for background processing
	h.JobQueue <- worker.ImportJob{FilePath: filename, ImportID: importID}

	return c.JSON(fiber.Map{
		"status":  "success",
		"data":    fiber.Map{"import_id": importID, "status": "queued"},
		"message": "File will be processed in the background",
	})
}

// GetImportStatusHandler retrieves the status of an import by its ID
func (h *ImportHandler) GetImportStatusHandler(c *fiber.Ctx) error {
	importID := c.Params("id")
	if importID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Import ID is required",
		})
	}

	status, err := h.ImportStatusRepo.FindByImportID(importID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Import status not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve import status: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   status,
	})
}
