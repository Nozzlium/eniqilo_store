package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/nozzlium/eniqilo_store/internal/constant"
	"github.com/nozzlium/eniqilo_store/internal/model"
	"github.com/nozzlium/eniqilo_store/internal/service"
)

type OrderHandler struct {
	OrderService *service.OrderService
}

func NewOrderHandler(
	orderService *service.OrderService,
) *OrderHandler {
	return &OrderHandler{
		OrderService: orderService,
	}
}

func (handlers *OrderHandler) Create(
	c *fiber.Ctx,
) error {
	var body model.OrderRequestBody
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
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
			return HandleError(
				c,
				ErrorResponse{
					message: "product quantity must be greater than 0",
					error:   constant.ErrBadInput,
					detail: fmt.Sprintf(
						"invalid quantity for product %s",
						product.ProductID,
					),
				},
			)
		}

		productId, err := uuid.Parse(
			product.ProductID,
		)
		if err != nil {
			return HandleError(
				c,
				ErrorResponse{
					message: "invalid product id",
					error:   constant.ErrBadInput,
					detail: fmt.Sprintf(
						"invalid product id: %s",
						product.ProductID,
					),
				},
			)
		}

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

	customerId, err := uuid.Parse(
		body.CustomerID,
	)
	if err != nil {
		return HandleError(
			c,
			ErrorResponse{
				message: "invalid customer id",
				error:   constant.ErrBadInput,
				detail: fmt.Sprintf(
					"invalid customer id: %s",
					body.CustomerID,
				),
			},
		)
	}

	res, err := handlers.OrderService.Create(
		c.Context(),
		model.Order{
			CustomerID:    customerId,
			PaymentAmount: body.Paid,
			Change:        body.Change,
			ProductOrders: productModels,
		},
	)
	if err != nil {
		return HandleError(
			c,
			ErrorResponse{
				error:   err,
				message: err.Error(),
				detail: fmt.Sprintf(
					"unable to create order %v",
					err,
				),
			},
		)
	}

	return c.JSON(fiber.Map{
		"message": "success",
		"data":    res.ToResponseBody(),
	})
}
