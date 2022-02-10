package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"statsv0/configs"
	"statsv0/models"
	"statsv0/rest/middlewares"
	"statsv0/tools/custerror"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	sorter "github.com/posener/order"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var greatestOrdersCollection *mongo.Collection = configs.GetCollection(configs.DB, "GreatestOrders")

func GuardarGreatestOrder(id string, month, year, articlesQ int) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	newGO := models.GreatestOrders{

		Id:              id,
		Month:           int(month),
		Year:            year,
		ArticleQuantity: articlesQ,
	}
	_, err := greatestOrdersCollection.InsertOne(ctx, newGO)
	if err != nil {
		println("Error inserting newGO")
		return
	}
	sortedOrders, err := SearchSortedOrders(year, month)
	if err != nil {
		println("Error searching sorted Orders")
		return
	}
	if len(sortedOrders) <= 10 {
		fmt.Println("no hay mas de 10 ordenes")
		return
	}
	fmt.Println("sorted order:", sortedOrders)
	lastOrder := sortedOrders[0]
	res, err := greatestOrdersCollection.DeleteOne(ctx, bson.D{{"id", bson.D{{"$eq", lastOrder.Id}}}})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("result: ", res, " id:", lastOrder.Id)
}

func toOrder(d primitive.D) models.GreatestOrders {

	var s models.GreatestOrders

	bsonBytes, _ := bson.Marshal(d)
	bson.Unmarshal(bsonBytes, &s)
	return s
}
func toOrden(d primitive.D) models.Orden {

	var s models.Orden

	bsonBytes, _ := bson.Marshal(d)
	bson.Unmarshal(bsonBytes, &s)
	return s
}

func GetGreatestOrdersSorted() gin.HandlerFunc {
	return func(c *gin.Context) {

		year, _ := strconv.Atoi(c.Param("year"))
		month, _ := strconv.Atoi(c.Param("month"))

		SortedOrders, err := SearchSortedOrders(year, month)
		if err != nil {
			c.Error(err)
			return
		}
		token, err := middlewares.GetHeaderToken(c)
		if err != nil {
			c.Error(err)
			return
		}
		fmt.Println(SortedOrders)
		var newOrders []models.GetGO
		for _, ord := range SortedOrders {
			fmt.Println("ord: ", ord)
			var newOrden models.Orden
			err := GetOrderData(ord.Id, token, &newOrden)
			if err != nil {
				c.Error(err)
			}
			respuesta := models.GetGO{

				Id:             newOrden.Id,
				Status:         newOrden.Status,
				Total:          newOrden.TotalPrice,
				Created:        newOrden.Created,
				TotalProductos: ord.ArticleQuantity,
			}
			newOrders = append(newOrders, respuesta)
		}
		var ordOrders = sorter.By(func(a, b models.GetGO) int { return b.TotalProductos - a.TotalProductos })
		ordOrders.Sort(newOrders)
		c.JSON(200, newOrders)

	}
}

func SearchSortedOrders(year, month int) ([]models.GreatestOrders, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	opts := options.Find()
	opts.SetSort(bson.D{{"articlequantity", 1}})
	defer cancel()
	sortCursor, err := greatestOrdersCollection.Find(ctx, bson.D{{"year", bson.D{{"$eq", year}}}, {"month", bson.D{{"$eq", month}}}}, opts)
	if err != nil {
		return nil, err

	}
	var greatestOrdersSorted []models.GreatestOrders
	if err = sortCursor.All(ctx, &greatestOrdersSorted); err != nil {
		return nil, err
	}

	return greatestOrdersSorted, nil
}

func GetOrderData(id string, token string, target *models.Orden) error {
	client := configs.Client()
	result, err := client.Get(id).Result()
	if err != nil {
		fmt.Println("no se encontro el registro en caché, buscando en el microservicio correspondiente...")
	} else {
		product, err := json.Marshal(result)
		if err != nil {
			return custerror.NewCustom(500, "no se pudo converir los datos de la orden en cache")
		}
		fmt.Println("retornando de cache:", product)
		uncuotedProduct, err := strconv.Unquote(string(product))
		if err != nil {
			fmt.Println(err)
		}
		return json.Unmarshal([]byte(uncuotedProduct), target)
	}
	req, err := http.NewRequest("GET", "http://localhost:3004/v1/orders/"+id, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "bearer "+token)
	response, err := http.DefaultClient.Do(req)
	fmt.Println("response:", response)
	if err != nil || response.StatusCode != 200 {
		return err
	}
	defer response.Body.Close()
	//guardamos el producto en cache una vez encontrado
	err = json.NewDecoder(response.Body).Decode(&target)
	if err != nil {
		return err
	}
	cacheProd, err := json.Marshal(&target)
	if err != nil {
		return err
	}
	fmt.Println("orden a cachear:", cacheProd)

	err = client.Set(id, cacheProd, 1*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}
