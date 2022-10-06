package middleware

import (
	"github.com/gin-gonic/gin"
	"gitlab.boquar.tech/galileosky/device/customer-administration/pkg/logger"
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
