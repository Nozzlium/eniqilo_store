package main

import (
	"log"

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

	userRepository := repository.NewUserRepository(db)

	userService := service.NewUserService(
		userRepository,
		cfg.JWTSecret,
		int(cfg.BCryptSalt),
	)

	err = handler.InitAuthHandler(
		fiberApp,
		userService,
	)
	if err != nil {
		log.Fatal(err)
	}

	err = fiberApp.Listen(":3000")
	if err != nil {
		log.Fatal(err)
	}
}
