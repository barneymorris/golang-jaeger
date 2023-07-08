package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/betelgeusexru/golang-jaeger/service-a/pkg/tracing"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	tracing.InitTracer()

	app := fiber.New()

	app.Use(otelfiber.Middleware())

	app.Get("/hello", func(ctx *fiber.Ctx) error {
		return ctx.JSON(map[string]string{"msg": "hello"})
	})

	app.Get("/external/bye", func(ctx *fiber.Ctx) error {
		request, err := http.NewRequestWithContext(ctx.Context(), "GET", "http://localhost:3001/bye", nil)
		if err != nil {
			return fmt.Errorf("create request error: %w", err)
		}

		client := http.Client{
			// Wrap the Transport with one that starts a span and injects the span context
			// into the outbound request headers.
			Transport: otelhttp.NewTransport(http.DefaultTransport),
			Timeout:   10 * time.Second,
		}

		r, err := client.Do(request)
		if err != nil {
			return fmt.Errorf("cannot make external call to service-b: %w", err)
		}

		response, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return fmt.Errorf("cannot read response body: %w", err)
		}

		type dto struct {
			Msg string `json:"msg"`
		}

		var mapped dto
		json.Unmarshal(response, &mapped)

		defer r.Body.Close()

		return ctx.JSON(mapped)
	})

	log.Fatal(app.Listen(":3000"))
}
