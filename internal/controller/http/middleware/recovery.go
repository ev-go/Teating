package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RecoveryMiddleware(c *gin.Context, recovered interface{}) {
	msg, ok := recovered.(string)
	if ok {
		msg = fmt.Sprintf("unexpected error: %s", msg)

		c.String(http.StatusInternalServerError, msg)
		c.Error(errors.New(msg))
	}

	c.AbortWithStatus(http.StatusInternalServerError)
}
