package handlers

import (
	"importerapi/internal/repositories"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type InvoiceItemHandler struct {
	InvoiceItemRepo repositories.InvoiceItemRepository
}

func NewInvoiceItemHandler(db *gorm.DB) *InvoiceItemHandler {
	return &InvoiceItemHandler{
		InvoiceItemRepo: repositories.NewInvoiceItemRepo(db),
	}
}

func (h *InvoiceItemHandler) GetInvoiceItemsHandler(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))

	filter := repositories.InvoiceItemFilter{
		InvoiceID:  c.Query("invoice_id"),
		ProductID:  c.Query("product_id"),
		CreditType: c.Query("credit_type"),
		OrderBy:    c.Query("orderBy"),
		OrderDir:   c.Query("order", "asc"),
		Page:       page,
		PageSize:   pageSize,
	}

	items, total, err := h.InvoiceItemRepo.FindAll(filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve invoices items",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":     "success",
		"data":       items,
		"total":      total,
		"page":       page,
		"pageSize":   pageSize,
		"totalPages": int(math.Ceil(float64(total) / float64(pageSize))),
	})
}
