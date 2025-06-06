package handlers

import (
	"fmt"
	"importerapi/internal/worker"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ImportHandler struct {
	JobQueue chan worker.ImportJob
}

func NewImportHandler(queue chan worker.ImportJob) *ImportHandler {
	return &ImportHandler{
		JobQueue: queue,
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

	// Add the job to the queue for background processing
	h.JobQueue <- worker.ImportJob{FilePath: filename}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "File will be processed in the background",
	})
}
