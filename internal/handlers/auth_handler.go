package handlers

import (
	"importerapi/internal/models"
	"importerapi/internal/repositories"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AuthHandler struct {
	UserRepo repositories.UserRepository
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{
		UserRepo: repositories.NewUserRepo(db),
	}
}

type requestBody struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

// RegisterHandler handles user registration.
func (h *AuthHandler) RegisterHandler(c *fiber.Ctx) error {
	var body requestBody

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error":   true,
			"message": "invalid request body",
		})
	}

	if body.Password != body.ConfirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": true, "message": "passwords doest match"})
	}

	u, err := models.NewUser(body.FirstName, body.LastName, body.Email, body.Password)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error":   true,
			"message": "error to create user",
		})
	}
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	if !models.IsValidEmail(u.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "invalid email format",
		})
	}

	err = h.UserRepo.Create(u)
	if err != nil {
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error":   true,
				"message": "email already exists",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "internal server error",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "user created successfully",
		"user": fiber.Map{
			"id":        u.ID,
			"firstName": u.FirstName,
			"lastName":  u.LastName,
			"email":     u.Email,
		},
	})
}
