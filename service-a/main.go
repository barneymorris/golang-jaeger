package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/betelgeusexru/golang-jaeger/service-a/pkg/tracing"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {

	tracing.InitTracer()

	r := gin.Default()
	r.Use(otelgin.Middleware("service-a"))

	r.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(200, map[string]string{"msg": "hello"})
	})

	r.GET("/external/bye", func(ctx *gin.Context) {

		request, err := http.NewRequestWithContext(ctx.Request.Context(), "GET", "http://localhost:3001/bye", nil)
		if err != nil {
			ctx.JSON(500, fmt.Errorf("create request error: %w", err))
		}

		client := http.Client{
			// Wrap the Transport with one that starts a span and injects the span context
			// into the outbound request headers.
			Transport: otelhttp.NewTransport(http.DefaultTransport),
			Timeout:   10 * time.Second,
		}

		r, err := client.Do(request)
		if err != nil {
			ctx.JSON(500, fmt.Errorf("cannot make external call to service-b: %w", err))
		}

		response, err := ioutil.ReadAll(r.Body)
		if err != nil {
			ctx.JSON(500, fmt.Errorf("cannot read response body: %w", err))
		}

		type dto struct {
			Msg string `json:"msg"`
		}

		var mapped dto
		json.Unmarshal(response, &mapped)

		defer r.Body.Close()

		ctx.JSON(200, mapped)

	})

	log.Fatal(r.Run(":3000"))
}
