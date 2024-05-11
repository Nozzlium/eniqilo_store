package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
		productModels = append(
			productModels,
			model.ProductOrder{
				ProductID: uuid.MustParse(
					product.ProductID,
				),
				Quantity: uint64(
					product.Quantity,
				),
			},
		)
	}

	if !body.IsValid() {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"message": "invalid body"})
	}

	res, err := handlers.OrderService.Create(
		c.UserContext(),
		model.Order{
			CustomerID: uuid.FromString(
				body.CustomerID,
			),
			PaymentAmount: body.Paid,
			Change:        body.Change,
			ProductOrders: productModels,
		},
	)
	if err != nil {
		if errors.Is(
			err,
			model.ErrNotFound,
		) {
			return c.Status(fiber.StatusNotFound).
				JSON(fiber.Map{
					"message": "not found",
				})
		}
		if errors.Is(
			err,
			model.ErrInsufficientFund,
		) {
			return c.Status(fiber.StatusBadRequest).
				JSON(fiber.Map{"message": err.Error()})
		}
		if errors.Is(
			err,
			model.ErrInsufficientStock,
		) {
			return c.Status(fiber.StatusBadRequest).
				JSON(fiber.Map{"message": err.Error()})
		}
		if errors.Is(
			err,
			model.ErrInvalidChange,
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
