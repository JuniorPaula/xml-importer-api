package handlers

import (
	"fmt"
	"importerapi/internal/util"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ImportHandler struct{}

func NewImportHandler() *ImportHandler {
	return &ImportHandler{}
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

	go func(path string) {
		f, err := os.Open(path)
		if err != nil {
			fmt.Printf("Error opening file: %s\n", err.Error())
			return
		}
		defer f.Close()

		fmt.Println("Initializing excel reader")
		records, err := util.ReadExcelFromReader(f)
		if err != nil {
			fmt.Println("Error reading excel file:", err)
			return
		}

		fmt.Printf("File %s processed successfully\n", path)
		fmt.Printf("Records %+v\n", records)

		// remove the file after processing
		if err := os.Remove(path); err != nil {
			fmt.Printf("Error removing file %s: %s\n", path, err.Error())
		} else {
			fmt.Printf("File %s removed successfully\n", path)
		}

	}(filename)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "File will be processed in the background",
	})
}
