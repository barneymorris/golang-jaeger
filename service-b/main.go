package main

import (
	"log"

	"github.com/betelgeusexru/golang-jaeger/service-b/pkg/tracing"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
)

func main() {
	tracing.InitTracer()

	app := fiber.New()

	app.Use(otelfiber.Middleware())

	app.Get("/bye", func(ctx *fiber.Ctx) error {
		return ctx.JSON(map[string]string{"msg": "bye"})
	})

	log.Fatal(app.Listen(":3001"))
}
