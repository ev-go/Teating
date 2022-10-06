package middleware

import (
	"github.com/gin-gonic/gin"
	"gitlab.boquar.tech/galileosky/pkg/acl"
	"net/http"
)

func ACLMiddleware(c *gin.Context) {
	//c.request.Header("A")
	hasAccess, err := acl.ACLRepo.HasAccess(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		c.Abort()
		return
	}

	if !hasAccess {
		c.JSON(http.StatusForbidden, "Operation not allowed")
		c.Abort()
		return
	} else {
		c.Next()
	}
	c.Next()
}
