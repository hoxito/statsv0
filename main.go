package main

import (
	"fmt"
	"statsv0/configs"
	"statsv0/rabbit"
	"statsv0/rest/middlewares"
	"statsv0/rest/routes"
	"statsv0/tools/env"
	"time"

	cors "github.com/itsjamie/gin-cors"

	"github.com/gin-gonic/gin"
)

func main() {
	rabbit.Init()
	rabbit.InitOrders()
	router := gin.Default()
	router.Use(middlewares.ErrorHandler)

	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type, Size",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	configs.ConnectDB()

	routes.PeakHourRoute(router)
	routes.GreatestOrdersRoute(router)
	routes.GreatestProductsRoute(router)
	routes.SellsPerDayRoute(router)
	router.Run(fmt.Sprintf(":%d", env.Get().Port))
}
