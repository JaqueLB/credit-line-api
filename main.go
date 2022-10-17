package main

import (
	"credit-line-api/src/api"
	"credit-line-api/src/api/entity"
	"credit-line-api/src/api/handler"
	"credit-line-api/src/db"

	"github.com/gin-gonic/gin"
)

func startHandlers() *api.Handlers {
	localStorage := make(map[int]*entity.CreditLineResponse)
	return &api.Handlers{
		CreditLineHandler: &handler.CreditLineHandler{
			Storage: &db.LocalStorage{
				Items: localStorage,
			},
		},
	}
}

func main() {
	engine := gin.Default()
	handlers := startHandlers()
	router := api.CreateRouter(handlers, engine)
	router.Run()
}
