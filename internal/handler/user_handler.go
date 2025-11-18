package handler

import (
	"strconv"

	"github.com/DMaryanskiy/go-idk/internal/domain"
	"github.com/DMaryanskiy/go-idk/internal/validator"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type UserHandler struct {
	service   domain.UserService
	validator *validator.Validator
	logger    *zap.Logger
}

func NewUserHandler(service domain.UserService, validator *validator.Validator, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		service:   service,
		validator: validator,
		logger:    logger,
	}
}

func (h *UserHandler) RegisterRoutes(router fiber.Router) {
	users := router.Group("/users")
	users.Post("/", h.CreateUser)
	users.Get("/", h.GetUsers)
	users.Get("/:id", h.GetUser)
	users.Put("/:id", h.UpdateUser)
	users.Delete("/:id", h.DeleteUser)
}

func (h *UserHandler) CreateUser(c fiber.Ctx) error {
	req := new(domain.CreateUserRequest)
	if err := c.Bind().JSON(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.validator.Validate(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	user, err := h.service.CreateUser(c.Context(), req)
	if err != nil {
		if err.Error() == "user with email already exists" {
			return fiber.NewError(fiber.StatusConflict, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create user")
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *UserHandler) GetUser(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	user, err := h.service.GetUser(c.Context(), id)
	if err != nil {
		if err.Error() == "user not found" {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get user")
	}

	return c.JSON(user)
}

func (h *UserHandler) GetUsers(c fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	response, err := h.service.GetUsers(c.Context(), limit, offset)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get users") 
	}

	return c.JSON(response)
}

func (h *UserHandler) UpdateUser(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	req := new(domain.UpdateUserRequest)
	if err := c.Bind().JSON(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.validator.Validate(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	user, err := h.service.UpdateUser(c.Context(), id, req)
	if err != nil {
		if err.Error() == "user not found" {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		if err.Error() == "email already in use" {
			return fiber.NewError(fiber.StatusConflict, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update user")
	}

	return c.JSON(user)
}

func (h *UserHandler) DeleteUser(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	if err := h.service.DeleteUser(c.Context(), id); err != nil {
		if err.Error() == "user not found" {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete user")
	}

	return c.SendStatus(fiber.StatusNoContent)
}
