package handler

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/nozzlium/eniqilo_store/internal/constant"
)

func HandleError(ctx *fiber.Ctx, err ErrorResponse) error {
	log.Println(err.detail)
	switch err {
	case constant.ErrNotFound:
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": err.message,
		})
	case constant.ErrConflict:
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": err.message,
		})
	case constant.ErrBadInput, constant.ErrInvalidBody, constant.ErrInsufficientFund, constant.ErrInvalidChange, constant.ErrInsufficientStock:
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.message,
		})
	default:
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "internal server error",
		})
	}
}

type ErrorResponse struct {
	error   error
	message string
	detail  string
}

func (e ErrorResponse) Error() string {
	return e.error.Error()
}
