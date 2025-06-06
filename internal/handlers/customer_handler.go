package handlers

import (
	"importerapi/internal/repositories"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CustomerHandler struct {
	CustomerRepo repositories.CustomerRepository
}

func NewCustomerHandler(db *gorm.DB) *CustomerHandler {
	return &CustomerHandler{
		CustomerRepo: repositories.NewCustomerRepo(db),
	}
}

func (h *CustomerHandler) GetCustomersHandler(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))

	filter := repositories.CustomerFilter{
		Name:     c.Query("name"),
		Country:  c.Query("country"),
		Page:     page,
		PageSize: pageSize,
	}

	customers, total, err := h.CustomerRepo.FindAll(filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve customers",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":     "success",
		"data":       customers,
		"total":      total,
		"page":       page,
		"pageSize":   pageSize,
		"totalPages": int(math.Ceil(float64(total) / float64(pageSize))),
	})
}
