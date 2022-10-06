package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.boquar.tech/galileosky/device/customer-administration/pkg/tracer"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func TracerMiddleware(c *gin.Context) {
	traceCtx, span := tracer.Start(
		c.Request.Context(),
		fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path),
		trace.WithSpanKind(trace.SpanKindServer),
	)
	defer tracer.End(span)

	c.Request = c.Request.Clone(traceCtx)
	c.Next()
	span.SetAttributes(attribute.Int("http.status.code", c.Writer.Status()))

	for _, err := range c.Errors.Errors() {
		tracer.Error(span, errors.New(err))
	}
}
