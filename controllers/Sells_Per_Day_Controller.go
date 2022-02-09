package controllers

import (
	"fmt"
	"statsv0/models"
	"statsv0/rest/middlewares"
	"statsv0/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetSellsPerDay() gin.HandlerFunc {
	return func(c *gin.Context) {
		var day *models.AggSellsPerDay
		day, err := services.SearchBestDay()
		if err != nil {
			c.Error(err)
		}
		fmt.Println("day:", day)

		c.JSON(200, day)
	}

}
func GetSellsPerDayYear() gin.HandlerFunc {
	return func(c *gin.Context) {
		year, _ := strconv.Atoi(c.Param("year"))
		var day *models.AggSellsPerDay
		day, err := services.SearchYearDay(year)
		if err != nil {
			c.Error(err)
		}
		fmt.Println("day:", day)

		c.JSON(200, day)
	}

}
func GetSellsPerDayMonth() gin.HandlerFunc {
	return func(c *gin.Context) {
		year, _ := strconv.Atoi(c.Param("year"))
		month, _ := strconv.Atoi(c.Param("month"))
		var day *models.AggSellsPerDay
		day, err := services.SearchMonthDay(year, month)
		if err != nil {
			c.Error(err)
		}
		fmt.Println("day:", day)

		c.JSON(200, day)
	}

}

func GetSellsPerDayProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		year, _ := strconv.Atoi(c.Param("year"))
		month, _ := strconv.Atoi(c.Param("month"))
		productId := c.Param("ProductId")
		var day *models.AggSellsPerDay
		day, err := services.SearchProductDay(year, month, productId)
		if err != nil {
			c.Error(err)
		}

		token, err := middlewares.GetHeaderToken(c)
		if err != nil {
			c.Error(err)
			return
		}
		var product models.Product
		services.GetProductData(productId, token, &product)
		fmt.Println("product:", product)
		fmt.Println("day:", day)
		var resp = models.SellsPerDayProduct{

			Weekday:     day.ID.Weekday,
			Month:       day.ID.Month,
			Year:        day.ID.Year,
			Quantity:    day.Quantity,
			Name:        product.Name,
			Description: product.Description,
			Image:       product.Image,
			Price:       product.Price,
			Stock:       product.Stock,
			Enabled:     product.Enabled,
		}
		c.JSON(200, resp)
	}

}
