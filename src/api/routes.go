package api

import (
	"credit-line-api/src/api/handler"
	"credit-line-api/src/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateRouter(r *gin.Engine) *gin.Engine {
	h := &handler.CreditLineHandler{
		Storage: &db.LocalStorage{},
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/credit-line/:clientID", h.Check)

	return r
}
