package controllers

import (
	"context"
	"fmt"
	"net/http"
	"statsv0/configs"
	"statsv0/models"
	"statsv0/responses"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var peakHourCollection *mongo.Collection = configs.GetCollection(configs.DB, "PeakHour")

func CreatePeakHour(c *gin.Context) {

	var peakHour models.PeakHour

	//valida el body
	if err := c.BindJSON(&peakHour); err != nil {
		c.Error(err)
		return
	}

	// Valida los campos requeridos usando la estructura estandar
	if validationErr := validate.Struct(&peakHour); validationErr != nil {
		c.Error(validationErr)
		return
	}

	newPH := models.PeakHour{

		Hour:          peakHour.Hour,
		Month:         peakHour.Month,
		Year:          peakHour.Year,
		OrderQuantity: peakHour.OrderQuantity,
	}

	result, err := peakHourCollection.InsertOne(c, newPH)
	if err != nil {
		c.Error(err)
		return
	}

	fmt.Println("creando...")
	c.JSON(http.StatusCreated, gin.H{"ID:": result.InsertedID, "hora": newPH.Hour, "mes": newPH.Month, "año": newPH.Year, "cantidad:": newPH.OrderQuantity})
}

func GetYearPeakHourSorted() gin.HandlerFunc {
	return func(c *gin.Context) {

		year, _ := strconv.Atoi(c.Param("year"))

		fmt.Println("year")
		fmt.Println(year)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		opts := options.Find()
		opts.SetSort(bson.D{{"orderquantity", -1}})
		opts.SetLimit(3)
		defer cancel()
		sortCursor, err := peakHourCollection.Find(ctx, bson.D{{"year", bson.D{{"$eq", year}}}}, opts)
		if err != nil {
			c.Error(err)
			return
		}
		fmt.Println(sortCursor.Current)
		var peakHourSorted []bson.M
		if err = sortCursor.All(ctx, &peakHourSorted); err != nil {
			c.Error(err)
			return
		}
		fmt.Println(peakHourSorted)

		c.JSON(200, peakHourSorted)
	}
}

func GetPeakHourSorted() gin.HandlerFunc {
	return func(c *gin.Context) {

		year, _ := strconv.Atoi(c.Param("year"))
		month, _ := strconv.Atoi(c.Param("month"))

		fmt.Println("year")
		fmt.Println(year)
		fmt.Println("month")
		fmt.Println(month)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		opts := options.Find()
		opts.SetSort(bson.D{{"orderquantity", -1}})
		opts.SetLimit(3)
		defer cancel()
		sortCursor, err := peakHourCollection.Find(ctx, bson.D{{"year", bson.D{{"$eq", year}}}, {"month", bson.D{{"$eq", month}}}}, opts)
		if err != nil {
			c.Error(err)
			return
		}
		fmt.Println(sortCursor.Current)
		var peakHourSorted []bson.M
		if err = sortCursor.All(ctx, &peakHourSorted); err != nil {
			c.Error(err)
			return
		}
		fmt.Println(peakHourSorted)

		c.JSON(200, peakHourSorted)
	}
}
func IncrementPeakHour() gin.HandlerFunc {
	return func(c *gin.Context) {
		year, _ := strconv.Atoi(c.Param("year"))
		month, _ := strconv.Atoi(c.Param("month"))
		hour, _ := strconv.Atoi(c.Param("hour"))
		s := fmt.Sprintf("fecha:%d/%d - %d", year, month, hour)
		println(s)

		err := IncPeakHour(year, month, hour)
		if err != nil {
			c.Error(err)
		}

		c.JSON(200, "Hora incrementada")
	}
}
func IncPeakHour(year, month, hour int) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	options := options.Update()
	options.SetUpsert(true)
	defer cancel()
	result, err := peakHourCollection.UpdateOne(ctx, bson.M{"year": year, "month": month, "hour": hour}, bson.D{{"$inc", bson.D{{"orderquantity", 1}}}}, options)
	if err != nil {
		return err
	}
	fmt.Println("incrementando la hora:", hour)
	fmt.Println("result:", result)
	return nil
}

func DeletePeakHour() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		year, _ := strconv.Atoi(c.Param("year"))
		month, _ := strconv.Atoi(c.Param("month"))
		hour, _ := strconv.Atoi(c.Param("hour"))
		s := fmt.Sprintf("fecha:%d/%d - %d", year, month, hour)
		println(s)
		defer cancel()

		result, err := peakHourCollection.DeleteOne(ctx, bson.M{"year": year, "month": month, "hour": hour})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.StdResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound,
				responses.StdResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "peakHour with specified datetime not found!"}},
			)
			return
		}

		c.JSON(http.StatusOK,
			responses.StdResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "User successfully deleted!"}},
		)
	}
}
