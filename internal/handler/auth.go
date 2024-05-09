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

type RegisterBody struct {
	PhoneNumber string `json:"email"`
	Name        string `json:"name"`
	Password    string `json:"password"`
}

type LoginBody struct {
	PhoneNumber string `json:"email"`
	Password    string `json:"password"`
}

func (handlers *AuthHandler) RegisterHandler(
	ctx *fiber.Ctx,
) error {
	var body RegisterBody
	err := ctx.BodyParser(&body)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{
				"message": err.Error(),
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
	var body LoginBody
	err := ctx.BodyParser(&body)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{
				"message": "unable to process body",
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
