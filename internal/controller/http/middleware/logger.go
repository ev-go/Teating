package middleware

import (
	"github.com/ev-go/Testing/pkg/logger"
	"github.com/gin-gonic/gin"
)

func LoggerMiddleware(c *gin.Context) {
	c.Next()

	l := logger.
		I.
		WithContext(c.Request.Context()).
		WithField("method", c.Request.Method).
		WithField("path", c.Request.URL.Path).
		WithField("ip", c.ClientIP()).
		WithField("proto", c.Request.Proto).
		WithField("userAgent", c.Request.UserAgent()).
		WithField("statusCode", c.Writer.Status()).
		WithField("responseSize", c.Writer.Size())

	for _, err := range c.Errors.Errors() {
		l.Error(err)
	}
}
