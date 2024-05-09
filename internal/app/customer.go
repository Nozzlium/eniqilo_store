package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/nozzlium/eniqilo_store/internal/model"
	"github.com/nozzlium/eniqilo_store/internal/repository"
	"github.com/nozzlium/eniqilo_store/internal/service"
)

type RegisterBody struct {
	PhoneNumber string `json:"phoneNumber"`
	Name        string `json:"name"`
}

func InitCustomer(
	app *fiber.App,
	conn *pgx.Conn,
) {
	customerRepository := repository.CustomerRepository{}
	customerService := service.CustomerService{
		CustomerRepository: &customerRepository,
		Conn:               conn,
	}
	customer := app.Group(
		"/v1/customer",
	)
	customer.Post(
		"/register",
		func(c *fiber.Ctx) error {
			var body RegisterBody
			err := c.BodyParser(&body)
			if err != nil {
				return err
			}

			customerData, err := customerService.Create(
				c.UserContext(),
				model.Customer{
					PhoneNumber: body.PhoneNumber,
					Name:        body.Name,
				},
			)
			if err != nil {
				return err
			}

			return c.JSON(fiber.Map{
				"message": "success",
				"data":    customerData,
			})
		},
	)
	customer.Get(
		"",
		func(c *fiber.Ctx) error {
			phoneNumber := c.Query(
				"phoneNumber",
				"",
			)
			name := c.Query("name", "")

			customerData, err := customerService.FindCustomers(
				c.UserContext(),
				model.Customer{
					PhoneNumber: phoneNumber,
					Name:        name,
				},
			)
			if err != nil {
				return err
			}

			return c.JSON(fiber.Map{
				"message": "success",
				"data":    customerData,
			})
		},
	)
}
