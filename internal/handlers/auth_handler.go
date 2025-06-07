package handlers

import (
	"importerapi/internal/models"
	"importerapi/internal/repositories"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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

type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginHandler handles user login and returns a JWT token.
func (h *AuthHandler) LoginHandler(c *fiber.Ctx) error {
	var body credentials

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error":   true,
			"message": "invalid request body",
		})
	}

	user, err := h.UserRepo.FindByEmail(body.Email)
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "unauthorized",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error":   true,
				"message": "internal server error",
			})
	}

	if !user.CheckPassword(body.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "unauthorized",
		})
	}

	claims := jwt.MapClaims{
		"sub":        user.ID,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error":   true,
			"message": "error signing token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":   false,
		"message": "Login successful",
		"data":    map[string]any{"user": user, "token": t},
	})
}
