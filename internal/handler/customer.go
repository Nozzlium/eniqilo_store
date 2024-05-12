package handler

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/nozzlium/eniqilo_store/internal/middleware"
	"github.com/nozzlium/eniqilo_store/internal/model"
	"github.com/nozzlium/eniqilo_store/internal/service"
)

type CustomerHandler struct {
	CustomerService *service.CustomerService
}

func NewCustomerHandler(
	customerService *service.CustomerService,
) *CustomerHandler {
	return &CustomerHandler{
		CustomerService: customerService,
	}
}

func InitCustomerHandler(
	app *fiber.App,
	customerService *service.CustomerService,
) error {
	if customerService == nil {
		return errors.New(
			"cannot init, customer service is nil",
		)
	}

	customerHandler := &CustomerHandler{
		CustomerService: customerService,
	}

	customer := app.Group("/customer")
	customer.Use(middleware.Protected())
	customer.Post(
		"/register",
		customerHandler.Register,
	)
	customer.Get(
		"",
		customerHandler.GetCustomers,
	)

	return nil
}

func (handler *CustomerHandler) Register(
	c *fiber.Ctx,
) error {
	var body model.CustomerRegisterBody
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"message": err.Error(),
			})
	}

	if !body.IsValid() {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"message": "request does not pass validation",
			})
	}

	customerData, err := handler.CustomerService.Create(
		c.Context(),
		model.Customer{
			PhoneNumber: body.PhoneNumber,
			Name:        body.Name,
		},
	)
	if err != nil {
		return HandleError(
			c,
			ErrorResponse{
				message: "error creating customer",
				error:   err,
				detail: fmt.Sprintf(
					"error creating customer: %v",
					err.Error(),
				),
			},
		)
	}

	return c.Status(fiber.StatusCreated).
		JSON(fiber.Map{
			"message": "success",
			"data":    customerData,
		})
}

func (handler *CustomerHandler) GetCustomers(
	c *fiber.Ctx,
) error {
	phoneNumber := c.Query(
		"phoneNumber",
		"",
	)
	name := c.Query("name", "")

	customerData, err := handler.CustomerService.FindCustomers(
		c.Context(),
		model.Customer{
			PhoneNumber: phoneNumber,
			Name:        name,
		},
	)
	if err != nil {
		return HandleError(
			c,
			ErrorResponse{
				message: "error get customer",
				error:   err,
				detail: fmt.Sprintf(
					"error get customer: %v",
					err.Error(),
				),
			},
		)
	}

	return c.JSON(fiber.Map{
		"message": "success",
		"data":    customerData,
	})
}
