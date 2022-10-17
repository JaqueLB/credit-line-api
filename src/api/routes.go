package api

import (
	"credit-line-api/src/api/handler"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	CreditLineHandler *handler.CreditLineHandler
}

func CreateRouter(h *Handlers, r *gin.Engine) *gin.Engine {
	r.POST("/credit-line/:clientID", h.CreditLineHandler.Check)

	return r
}
