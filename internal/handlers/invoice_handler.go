package handlers

import (
	"importerapi/internal/repositories"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type InvoiceHandler struct {
	InvoiceRepo repositories.InvoiceRepository
}

func NewInvoiceHandler(db *gorm.DB) *InvoiceHandler {
	return &InvoiceHandler{
		InvoiceRepo: repositories.NewInvoiceRepo(db),
	}
}

func (h *InvoiceHandler) GetInvoicesHandler(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))
	month, _ := strconv.Atoi(c.Query("month"))
	year, _ := strconv.Atoi(c.Query("year"))

	filter := repositories.InvoiceFilter{
		CustomerID: c.Query("customer_id"),
		Month:      month,
		Year:       year,
		Page:       page,
		PageSize:   pageSize,
	}

	invoices, total, err := h.InvoiceRepo.FindAll(filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve invoices",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":     "success",
		"data":       invoices,
		"total":      total,
		"page":       page,
		"pageSize":   pageSize,
		"totalPages": int(math.Ceil(float64(total) / float64(pageSize))),
	})
}
