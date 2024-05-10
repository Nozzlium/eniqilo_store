package main

import (
	"log"
	"regexp"

	"github.com/bytedance/sonic"
	"github.com/caarlos0/env/v11"
	"github.com/gofiber/fiber/v2"
	"github.com/nozzlium/eniqilo_store/internal/client"
	"github.com/nozzlium/eniqilo_store/internal/config"
	"github.com/nozzlium/eniqilo_store/internal/handler"
	"github.com/nozzlium/eniqilo_store/internal/repository"
	"github.com/nozzlium/eniqilo_store/internal/service"
)

func main() {
	fiberApp := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
		Prefork:     true,
	})

	var cfg config.Config
	opts := env.Options{
		TagName: "json",
	}
	if err := env.ParseWithOptions(&cfg, opts); err != nil {
		log.Fatalf("%+v\n", err)
	}

	db, err := client.InitDB(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	userRepository := repository.NewUserRepository(
		db,
	)
	customerRepository, err := repository.NewCustomerRepository(
		db,
	)
	if err != nil {
		log.Fatal(err)
	}

	userService := service.NewUserService(
		userRepository,
		cfg.JWTSecret,
		int(cfg.BCryptSalt),
	)
	customerService, err := service.NewCustomerService(
		customerRepository,
		regexp.MustCompile(
			"^[+]{1}[0-9]{10,15}$",
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = handler.InitAuthHandler(
		fiberApp,
		userService,
	)
	if err != nil {
		log.Fatal(err)
	}
	err = handler.InitCustomerHandler(
		fiberApp,
		customerService,
	)
	if err != nil {
		log.Fatal(err)
	}

	err = fiberApp.Listen(":3000")
	if err != nil {
		log.Fatal(err)
	}
}
