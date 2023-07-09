package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/betelgeusexru/golang-jaeger/service-a/pkg/tracing"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	log.Info("Serivce a: initing the service-a")
	tracing.InitTracer()

	r := gin.Default()
	r.Use(otelgin.Middleware("service-a"))

	r.GET("/error", func(ctx *gin.Context) {
		span := trace.SpanFromContext(ctx.Request.Context())
		traceID := span.SpanContext().TraceID().String()

		log.WithFields(log.Fields{
			"trace-id": traceID,
		}).Info("Service a: hit GET /error endpoint")

		log.WithFields(log.Fields{
			"trace-id": traceID,
		}).Error("Service a: GET /error endpoint return a error")

		ctx.JSON(500, map[string]string{"msg": "error"})
	})

	r.GET("/hello", func(ctx *gin.Context) {
		span := trace.SpanFromContext(ctx.Request.Context())
		traceID := span.SpanContext().TraceID().String()

		log.WithFields(log.Fields{
			"trace-id": traceID,
		}).Info("Service a: hit GET /hello endpoint")

		log.WithFields(log.Fields{
			"trace-id": traceID,
			"response": map[string]string{"msg": "hello"},
		}).Info("Service a: endpoint GET /hello successfully done")
		ctx.JSON(200, map[string]string{"msg": "hello"})
	})

	r.GET("/external/bye", func(ctx *gin.Context) {
		span := trace.SpanFromContext(ctx.Request.Context())
		traceID := span.SpanContext().TraceID().String()

		log.WithFields(log.Fields{
			"trace-id": traceID,
		}).Info("Service a: hit GET /external/bye endpoint")

		request, err := http.NewRequestWithContext(ctx.Request.Context(), "GET", "http://localhost:3001/bye", nil)
		if err != nil {
			log.WithFields(log.Fields{
				"trace-id": traceID,
			}).Errorf("Service a: create request error: %s", err)
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
			log.WithFields(log.Fields{
				"trace-id": traceID,
			}).Errorf("Service a: cannot make external call to service-b: %s", err)
			ctx.JSON(500, fmt.Errorf("cannot make external call to service-b: %w", err))
		}

		response, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.WithFields(log.Fields{
				"trace-id": traceID,
			}).Errorf("Service a: cannot read response body: %s", err)
			ctx.JSON(500, fmt.Errorf("cannot read response body: %w", err))
		}

		type dto struct {
			Msg string `json:"msg"`
		}

		var mapped dto
		json.Unmarshal(response, &mapped)

		defer r.Body.Close()

		log.WithFields(log.Fields{
			"trace-id": traceID,
			"response": mapped,
		}).Info("Service a: endpoint GET /external/bye successfully done")
		ctx.JSON(200, mapped)

	})

	log.Info("Serivce a: starting service on port 3000...")
	err := r.Run(":3000")
	if err != nil {
		log.Fatal("Service a: cannot run service-a on port 3000")
	}
}
