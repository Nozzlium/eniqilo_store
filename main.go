package main

import (
	"log"

	"github.com/bytedance/sonic"
	"github.com/caarlos0/env/v11"
	"github.com/gofiber/fiber/v2"
	"github.com/nozzlium/eniqilo_store/internal/client"
	"github.com/nozzlium/eniqilo_store/internal/config"
	"github.com/nozzlium/eniqilo_store/internal/handler"
	"github.com/nozzlium/eniqilo_store/internal/middleware"
	"github.com/nozzlium/eniqilo_store/internal/repository"
	"github.com/nozzlium/eniqilo_store/internal/service"
)

func main() {
	fiberApp := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
		Prefork:     true,
	})

	err := setupApp(fiberApp)
	if err != nil {
		log.Fatal(err)
	}

	err = fiberApp.Listen(":3000")
	if err != nil {
		log.Fatal(err)
	}
}

func setupApp(app *fiber.App) error {
	var cfg config.Config
	opts := env.Options{
		TagName: "json",
	}
	if err := env.ParseWithOptions(&cfg, opts); err != nil {
		log.Fatalf("%+v\n", err)
		return err
	}

	db, err := client.InitDB(cfg.DB)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// initiate repositories
	userRepository := repository.NewUserRepository(
		db,
	)
	productRepository := repository.NewProductRepository(
		db,
	)
	customerRepository, _ := repository.NewCustomerRepository(
		db,
	)

	// initiate services
	userService := service.NewUserService(
		userRepository,
		cfg.JWTSecret,
		int(cfg.BCryptSalt),
	)

	productService := service.NewProductService(
		productRepository,
	)
	customerService, _ := service.NewCustomerService(
		customerRepository,
	)

	// initiate handlers
	authHandler := handler.NewAuthHandler(
		userService,
	)
	productHandler := handler.NewProductHandler(
		productService,
	)
	customerHandler := handler.NewCustomerHandler(
		customerService,
	)

	v1 := app.Group("/v1")
	auth := v1.Group("/staff")
	auth.Post(
		"/register",
		authHandler.RegisterHandler,
	)

	auth.Post(
		"/login",
		authHandler.Login,
	)

	product := v1.Group("/product")
	product.Get(
		"/customer",
		productHandler.SearchForCustomer,
	)

	// protected routes (require authentication)
	protectedProduct := product.Use(middleware.Protected()).
		Use(middleware.SetEmailAndUserID())
	protectedProduct.Get(
		"",
		productHandler.Search,
	)
	protectedProduct.Post(
		"",
		productHandler.Create,
	)
	protectedProduct.Put(
		"/:id",
		productHandler.Update,
	)
	protectedProduct.Delete(
		"/:id",
		productHandler.Delete,
	)

	customer := app.Group(
		"/v1/customer",
	)
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
