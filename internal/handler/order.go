package handler

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/nozzlium/eniqilo_store/internal/constant"
	"github.com/nozzlium/eniqilo_store/internal/model"
	"github.com/nozzlium/eniqilo_store/internal/service"
)

type OrderHandlers struct {
	OrderService service.OrderService
}

func (handlers *OrderHandlers) Create(
	c *fiber.Ctx,
) error {
	var body model.OrderRequestBody
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{
				"message": err.Error(),
			})
	}

	productModels := make(
		[]model.ProductOrder,
		0,
		len(body.ProductDetails),
	)
	for _, product := range body.ProductDetails {
		if !product.IsValid() {
			return HandleError(c, ErrorResponse{
				message: "product quantity must be greater than 0",
				error:   constant.ErrBadInput,
				detail:  fmt.Sprintf("invalid quantity for product %s", product.ProductID),
			})
		}

		productId := uuid.MustParse(product.ProductID)
		productModels = append(
			productModels,
			model.ProductOrder{
				ProductID: productId,
				Quantity:  product.Quantity,
			},
		)
	}

	if !body.IsValid() {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"message": "invalid body"})
	}

	customerId := uuid.MustParse(body.CustomerID)
	res, err := handlers.OrderService.Create(
		c.UserContext(),
		model.Order{
			CustomerID:    customerId,
			PaymentAmount: body.Paid,
			Change:        body.Change,
			ProductOrders: productModels,
		},
	)
	if err != nil {
		if errors.Is(
			err,
			constant.ErrNotFound,
		) {
			return c.Status(fiber.StatusNotFound).
				JSON(fiber.Map{
					"message": "not found",
				})
		}
		if errors.Is(
			err,
			constant.ErrInsufficientFund,
		) {
			return c.Status(fiber.StatusBadRequest).
				JSON(fiber.Map{"message": err.Error()})
		}
		if errors.Is(
			err,
			constant.ErrInsufficientStock,
		) {
			return c.Status(fiber.StatusBadRequest).
				JSON(fiber.Map{"message": err.Error()})
		}
		if errors.Is(
			err,
			constant.ErrInvalidChange,
		) {
			return c.Status(fiber.StatusBadRequest).
				JSON(fiber.Map{"message": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "success",
		"data":    res.ToResponseBody(),
	})
}
