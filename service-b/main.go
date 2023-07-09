package main

import (
	"github.com/betelgeusexru/golang-jaeger/service-b/pkg/tracing"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/trace"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Serivce b: initing the service-b")
	tracing.InitTracer()

	r := gin.Default()
	r.Use(otelgin.Middleware("service-b"))

	r.GET("/bye", func(ctx *gin.Context) {
		span := trace.SpanFromContext(ctx.Request.Context())
		traceID := span.SpanContext().TraceID().String()

		log.WithFields(log.Fields{
			"trace-id": traceID,
		}).Info("Service b: hit GET /bye endpoint")

		log.WithFields(log.Fields{
			"trace-id": traceID,
			"response": map[string]string{"msg": "bye"},
		}).Info("Service b: endpoint GET /bye successfully done")

		ctx.JSON(200, map[string]string{"msg": "bye"})
	})

	log.Info("Serivce n: starting service on port 3001...")
	err := r.Run(":3001")
	if err != nil {
		log.Fatal("Service b: cannot run service-b on port 3001")
	}
}
