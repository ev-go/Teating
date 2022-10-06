package v1

import (
	"github.com/gin-gonic/gin"
)

type response struct {
	Error string `json:"error" example:"message"`
}

func errorResponse(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, response{msg})
}

func respondWithError(c *gin.Context, code int, err error) {
	c.AbortWithStatusJSON(code, errorResp{err.Error()})
	c.Error(err)
}

type errorResp struct {
	Error string `json:"error"`
}
