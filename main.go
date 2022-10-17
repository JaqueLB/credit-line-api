package main

import (
	"credit-line-api/src/api"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	router := api.CreateRouter(engine)
	// TODO: rejectedCount := 0
	router.Run()
}
