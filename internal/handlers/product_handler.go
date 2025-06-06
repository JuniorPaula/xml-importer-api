package handlers

import (
	"importerapi/internal/repositories"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ProductHandler struct {
	ProductRepo repositories.ProductRepository
}

func NewProductHandler(db *gorm.DB) *ProductHandler {
	return &ProductHandler{
		ProductRepo: repositories.NewProductRepo(db),
	}
}

func (h *ProductHandler) GetProducts(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))

	filter := repositories.ProductFilter{
		Name:        c.Query("name"),
		PublisherID: c.Query("publisher_id"),
		Page:        page,
		PageSize:    pageSize,
	}

	products, total, err := h.ProductRepo.FindAll(filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve products",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":     "success",
		"data":       products,
		"total":      total,
		"page":       page,
		"pageSize":   pageSize,
		"totalPages": int(math.Ceil(float64(total) / float64(pageSize))),
	})
}
