package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/api/get-stock-list", GetStockList)
	router.GET("/api/goland", Goland)
	router.Run(":8081")
}
