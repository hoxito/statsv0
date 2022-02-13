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
	"statsv0/tools/env"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	sorter "github.com/posener/order"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var greatestOrdersCollection *mongo.Collection = configs.GetCollection(configs.DB, "GreatestOrders")

//Se ejecuta al realizarse una ordenplaced desde el microservicio de ordenes. Este evento se comunica por rabbitMQ y el topic consumer llama a esta funcion.

//Guarda los datos de la orden ingresados (id, mes año y la cantidad de articulos vendidos)
//en la base de datos de mongoDB. Luego de guardarla, si hay mas de 10 ordenes en la base de datos,
// busca la orden con menor cantidad de articulos y la elimina asegurando que siempre hayan 10 ordenes como maximo en la base de datos.
func GuardarGreatestOrder(id string, month, year, articlesQ int) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	newGO := models.GreatestOrders{

		Id:              id,
		Month:           int(month),
		Year:            year,
		ArticleQuantity: articlesQ,
	}

	sortedOrders, err := SearchSortedOrders(year, month)
	if err != nil {
		println("Error searching sorted Orders")
		return
	}

	if sortedOrders[0].ArticleQuantity > articlesQ && len(sortedOrders) >= 10 {
		// La orden ingresada no posee mas articulos vendidos que la ultima orden, por lo tanto, no entra en el ranking de las 10 mejores ordenes.
		return
	}
	_, err = greatestOrdersCollection.InsertOne(ctx, newGO)

	if err != nil {
		println("Error inserting newGO")
		return
	}

	sortedOrders, err = SearchSortedOrders(year, month) //se chequea la lista de ordenes una vez mas en caso de que durante la insercion se haya producido una concurrencia
	if err != nil {
		println("Error searching sorted Orders again")
		return
	}
	if len(sortedOrders) <= 10 {
		//Si no hay mas de 10 ordenes retorna
		return
	}
	//caso contrario, se elimina la ultima orden ordenada por cantidad de articulos vendidos
	_, err = greatestOrdersCollection.DeleteOne(ctx, bson.D{{"id", bson.D{{"$eq", sortedOrders[0].Id}}}})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Orden borrada:", sortedOrders[0].Id, " con ", sortedOrders[0].ArticleQuantity, " articulos vendidos")
}

// Trae las 10 mejores ordenes ordenadas segun la cantidad de articulos comprados
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

// Busca y recupera de la base de datos de mongo, las ordenes existentes ordenadas segun el campo "articlequantity" de un año y mes especificos
// Toma como argumentos el mes {month} y el año {year} de las ordenes como numeros enteros y devuelve una coleccion de "GreatestOrders" y un error
//En caso de no encontrar ordenes, retorna la coleccion vacía.
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

//Obtiene los datos de la orden con el id ingresado en los argumentos de la funcion
//Utiliza un token de autenticación para buscar las ordenes correspondientes al usuario autenticado.
//Para la obtencion de los datos, primero busca los datos de la orden por id en la base de datos redis utilizandola como caché.
//Si esta orden no es encontrada, luego busca en el microservicio "orders"
//El resultado es guardado en la variable pasada como argumento "target" de tipo "Orden".
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
	req, err := http.NewRequest("GET", env.Get().OrdersURL+"/v1/orders/"+id, nil)
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
