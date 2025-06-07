package handlers

import (
	"errors"
	"importerapi/internal/repositories"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserHandler struct {
	UserRepo repositories.UserRepository
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{
		UserRepo: repositories.NewUserRepo(db),
	}
}

// GetProfileHanlder retrieves the profile of the currently authenticated user.
func (h *UserHandler) GetProfileHanlder(c *fiber.Ctx) error {
	userID, err := getUserIdFromCtx(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "unauthorized",
		})
	}

	user, err := h.UserRepo.FindByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "user not found",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   user,
	})
}

// getUserIdFromCtx retrieves the user ID from the context.
// It expects the user information to be stored in the context as a fiber.Map.
// If the user is not found or the format is incorrect, it returns an error.
func getUserIdFromCtx(c *fiber.Ctx) (int, error) {
	userCtx := c.Locals("user")
	if userCtx == "" {
		return 0, errors.New("anauthoried")
	}

	userMap, ok := userCtx.(fiber.Map)
	if !ok {
		return 0, errors.New("invalid user format")
	}
	idInterface := userMap["id"]

	var userID int
	switch v := idInterface.(type) {
	case float64:
		userID = int(v)
	case int:
		userID = v
	case string:
		parsedID, err := strconv.Atoi(v)
		if err != nil {
			return 0, errors.New("invalid user ID format")
		}
		userID = parsedID
	default:
		return 0, errors.New("unexpected user ID type")
	}
	return userID, nil
}
