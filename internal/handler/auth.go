package handler

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/nozzlium/eniqilo_store/internal/model"
	"github.com/nozzlium/eniqilo_store/internal/service"
)

type AuthHandler struct {
	UserService *service.UserService
}

func NewAuthHandler(userService *service.UserService) *AuthHandler {
	return &AuthHandler{UserService: userService}
}

func InitAuthHandler(
	app *fiber.App,
	userService *service.UserService,
) error {
	if userService == nil {
		return errors.New(
			"cannot init, userService is nil",
		)
	}

	authHandler := AuthHandler{
		UserService: userService,
	}

	auth := app.Group("")
	auth.Post(
		"/register",
		authHandler.RegisterHandler,
	)

	auth.Post(
		"/login",
		authHandler.Login,
	)

	return nil
}

func (handlers *AuthHandler) RegisterHandler(
	ctx *fiber.Ctx,
) error {
	var body model.RegisterBody
	err := ctx.BodyParser(&body)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{
				"message": err.Error(),
			})
	}

	if !body.IsValid() {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"message": "invalid body",
			})
	}

	data, err := handlers.UserService.Register(
		ctx.Context(),
		model.User{
			PhoneNumber: body.PhoneNumber,
			Name:        body.Name,
			Password:    body.Password,
		},
	)
	if err != nil {
		if errors.Is(
			err,
			model.ErrConflict,
		) {
			return ctx.Status(fiber.StatusConflict).
				JSON(fiber.Map{
					"message": err.Error(),
				})
		}

		return ctx.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{
				"message": err.Error(),
			})
	}

	return ctx.JSON(fiber.Map{
		"message": "User registered successfully",
		"data":    data,
	})
}

func (handlers *AuthHandler) Login(
	ctx *fiber.Ctx,
) error {
	var body model.LoginBody
	err := ctx.BodyParser(&body)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{
				"message": "unable to process body",
			})
	}

	if !body.IsValid() {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"message": "invalid body",
			})
	}

	data, err := handlers.UserService.Login(
		ctx.Context(),
		model.User{
			PhoneNumber: body.PhoneNumber,
			Password:    body.Password,
		},
	)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{
				"message": err.Error(),
			})
	}

	return ctx.JSON(fiber.Map{
		"message": "User logged in successfully",
		"data":    data,
	})
}
