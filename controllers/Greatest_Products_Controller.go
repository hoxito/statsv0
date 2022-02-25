package controllers

import (
	"context"
	"fmt"
	"statsv0/configs"
	"statsv0/models"
	"statsv0/rest/middlewares"
	"statsv0/services"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	sorter "github.com/posener/order"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var greatestProductsCollection *mongo.Collection = configs.GetCollection(configs.DB, "GreatestProducts")

// guarda el producto con sus ventas en el año y mes en que se vendió. Si ya se había vendido un producto en el año y mes actuales, entonces lo actualiza sumando
// las ventas al producto ya guardado.
func GuardarGreatestProduct(art models.Article, month int, year int) {

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	opts := options.Update()
	opts.SetUpsert(true)
	result, err := greatestProductsCollection.UpdateOne(ctx, bson.M{"year": year, "month": int(month), "productid": art.ArticleId}, bson.D{{"$inc", bson.D{{"sells", art.Quantity}}}}, opts)
	if err != nil {
		println("Error inserting newGO")
		return
	}
	if result.MatchedCount == 0 {
		fmt.Println("Guardando nuevo producto vendido id: ", art.ArticleId)
	} else {
		fmt.Println("Sumando ", art.Quantity, "ventas al producto ", art.ArticleId)
	}
}

func GetProducts(what string) gin.HandlerFunc {
	return func(c *gin.Context) {
		year, _ := strconv.Atoi(c.Param("year"))
		month, _ := strconv.Atoi(c.Param("month"))
		best10, err := Get10(what, year, month) //Obtiene los 10 mejores o peores productos de la base de datos
		if err != nil {
			c.Error(err)
		}

		token, err := middlewares.GetHeaderToken(c)
		if err != nil {
			c.Error(err)
			return
		}

		var BestComplete []models.GPProduct //lista de productos con los datos completos de cada uno segun el endpoint consultado

		for _, product := range best10 {
			var producto models.Product
			services.GetProductData(product.ID.ProductId, token, &producto)
			fmt.Println("product:", product)
			var resp = models.GPProduct{

				ProductId:   product.ID.ProductId,
				Month:       product.ID.Month,
				Year:        product.ID.Year,
				TotalSells:  product.Sells,
				Name:        producto.Name,
				Description: producto.Description,
				Image:       producto.Image,
				Price:       producto.Price,
				Stock:       producto.Stock,
				Enabled:     producto.Enabled,
			}
			BestComplete = append(BestComplete, resp)
		}
		fmt.Println("best10:", best10)

		var ordProds = sorter.By(func(a, b models.GPProduct) int {
			if what == "best" {
				return b.TotalSells - a.TotalSells
			} else if what == "worst" {
				return a.TotalSells - b.TotalSells
			}
			return 0
		})
		ordProds.Sort(BestComplete)
		c.JSON(200, BestComplete)
	}

}

// Retorna los mejores o peores 10 productos segun ventas.
// El argumento "what" define si se retornaran los mejores "best" o peores "worst" 10 productos
func Get10(what string, year, month int) ([]models.AggGreatestProducts, error) {
	var sort int
	if what == "best" {
		sort = -1
	} else if what == "worst" {
		sort = 1
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	pipeline := []bson.M{
		{
			"$group": bson.M{
				"_id":           bson.M{"productid": "$productid", "month": "$month", "year": "$year"},
				"ventasTotales": bson.M{"$sum": "$sells"},
			},
		}, {
			"$sort": bson.M{"ventasTotales": sort},
		}, {
			"$limit": 10,
		},
	}
	sortCursor, err := greatestProductsCollection.Aggregate(ctx, pipeline)
	if err != nil {
		fmt.Println("error al bsucar en bd")
		return nil, err

	}
	var best10 []models.AggGreatestProducts

	if err = sortCursor.All(ctx, &best10); err != nil {
		fmt.Println("error al asignar el best10")
		return nil, err
	}

	return best10, nil
}
