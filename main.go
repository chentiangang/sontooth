package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/son/get-stock-list", GetStockList)
	router.GET("/son/goland", Goland)
	router.GET("/son/suggest", Suggest)
	router.GET("/ops/api/getHostNameByInstanceID", GetHostNameByInstanceID)
	router.Run(":8082")
}
