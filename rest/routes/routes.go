package routes

import (
	"statsv0/controllers"
	"statsv0/rest/middlewares"

	"github.com/gin-gonic/gin"
)

func PeakHourRoute(router *gin.Engine) {
	router.POST("/v1/stats/peakhour", middlewares.ValidateAuthentication, controllers.CreatePeakHour)
	router.GET("/v1/stats/peakhour/:year", controllers.GetYearPeakHourSorted())
	router.GET("/v1/stats/peakhour/:year/:month", controllers.GetPeakHourSorted())
	router.PUT("/v1/stats/peakhour/increment/:year/:month/:hour", middlewares.ValidateAuthentication, controllers.IncrementPeakHour())
	router.DELETE("/v1/stats/peakhour/:year/:month/:hour", middlewares.ValidateAuthentication, controllers.DeletePeakHour())

}

func GreatestOrdersRoute(router *gin.Engine) {
	router.GET("/v1/stats/greatestorders/:year/:month", controllers.GetGreatestOrdersSorted())
}

func GreatestProductsRoute(router *gin.Engine) {
	router.GET("/v1/stats/greatestproducts/:year/:month", controllers.GetProducts("best"))
	router.GET("/v1/stats/saddestproducts/:year/:month", controllers.GetProducts("worst"))
}

func SellsPerDayRoute(router *gin.Engine) {
	router.GET("/v1/stats/SellsPerDay/:year/:month/:ProductId", controllers.GetSellsPerDayProduct())
	router.GET("/v1/stats/SellsPerDay/:year/:month", controllers.GetSellsPerDayMonth())
	router.GET("/v1/stats/SellsPerDay/:year", controllers.GetSellsPerDayYear())
	router.GET("/v1/stats/SellsPerDay", controllers.GetSellsPerDay())
}
