package main

import (
	"log"

	"github.com/betelgeusexru/golang-jaeger/service-b/pkg/tracing"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
	tracing.InitTracer()

	r := gin.Default()
	r.Use(otelgin.Middleware("service-a"))

	r.GET("/bye", func(ctx *gin.Context) {
		ctx.JSON(200, map[string]string{"msg": "bye"})
	})

	log.Fatal(r.Run(":3001"))
}
