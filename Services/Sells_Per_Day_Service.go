package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"statsv0/configs"
	"statsv0/models"

	"statsv0/tools/custerror"
	"statsv0/tools/env"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var SellsPerDayCollection *mongo.Collection = configs.GetCollection(configs.DB, "SellsPerDay")



func GuardarSellsPerDay(art models.Article, month int, year int, weekday int) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	newSPD := models.SellsPerDay{

		ProductId: art.ArticleId,
		Month:     int(month),
		Year:      year,
		Weekday:   weekday,
		Quantity:  art.Quantity,
	}

	_, err := SellsPerDayCollection.InsertOne(ctx, newSPD)
	if err != nil {
		println("Error inserting newGO")
		return
	}
}
//encuentra dia de la semana con mas ventas de todos los articulos de todos los tiempos
func SearchBestDay() (*models.AggSellsPerDay, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pipeline := []bson.M{
		{
			"$group": bson.M{
				"_id":             bson.M{"weekday": "$weekday"},
				"articlequantity": bson.M{"$sum": "$quantity"},
			},
		},
		{
			"$sort": bson.M{"articlequantity": -1},
		},
	}
	sortCursor, err := SellsPerDayCollection.Aggregate(ctx, pipeline)
	if err != nil {
		fmt.Println("error al bsucar en bd")
		return nil, err

	}
	var productDay []models.AggSellsPerDay

	if err = sortCursor.All(ctx, &productDay); err != nil {
		fmt.Println("error al asignar el productDay")
		return nil, err
	}

	return &productDay[0], nil
}

//encuentra dia de la semana con mas ventas de todos los articulos del año especificado
func SearchYearDay(year int) (*models.AggSellsPerDay, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"year": year,
			},
		},
		{
			"$group": bson.M{
				"_id":             bson.M{"weekday": "$weekday", "year": "$year"},
				"articlequantity": bson.M{"$sum": "$quantity"},
			},
		},
		{
			"$sort": bson.M{"articlequantity": -1},
		},
	}
	sortCursor, err := SellsPerDayCollection.Aggregate(ctx, pipeline)
	if err != nil {
		fmt.Println("error al bsucar en bd")
		return nil, err

	}
	var productDay []models.AggSellsPerDay

	if err = sortCursor.All(ctx, &productDay); err != nil {
		fmt.Println("error al asignar el productDay")
		return nil, err
	}

	return &productDay[0], nil
}

//encuentra dia de la semana con mas ventas de todos los articulos del mes y año especificados
func SearchMonthDay(year, month int) (*models.AggSellsPerDay, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"year":  year,
				"month": month,
			},
		},
		{
			"$group": bson.M{
				"_id":             bson.M{"weekday": "$weekday", "month": "$month", "year": "$year"},
				"articlequantity": bson.M{"$sum": "$quantity"},
			},
		},
		{
			"$sort": bson.M{"articlequantity": -1},
		},
	}
	sortCursor, err := SellsPerDayCollection.Aggregate(ctx, pipeline)
	if err != nil {
		fmt.Println("error al bsucar en bd")
		return nil, err

	}
	var productDay []models.AggSellsPerDay

	if err = sortCursor.All(ctx, &productDay); err != nil {
		fmt.Println("error al asignar el productDay")
		return nil, err
	}

	return &productDay[0], nil
}

//encuentra dia de la semana con mas ventas de un articulo en especifico en un mes y año especifico
func SearchProductDay(year, month int, productId string) (*models.AggSellsPerDay, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"year":      year,
				"month":     month,
				"productid": productId,
			},
		},
		{
			"$group": bson.M{
				"_id":             bson.M{"weekday": "$weekday", "month": "$month", "year": "$year"},
				"articlequantity": bson.M{"$sum": "$quantity"},
			},
		},
		{
			"$sort": bson.M{"articlequantity": -1},
		},
	}
	sortCursor, err := SellsPerDayCollection.Aggregate(ctx, pipeline)
	if err != nil {
		fmt.Println("error al bsucar en bd")
		return nil, err

	}
	var productDay []models.AggSellsPerDay

	if err = sortCursor.All(ctx, &productDay); err != nil {
		fmt.Println("error al asignar el productDay")
		return nil, err
	}
	if len(productDay) == 0 {
		return nil, custerror.NewCustom(500, "No se encontró el registro con el id ingresado")
	}
	return &productDay[0], nil
}


//Funcion que obtiene de la cache o en caso de no ser posible, del microservicio catalog, los datos de un cierto producto
// a partir de su ID.
func GetProductData(id string, token string, target *models.Product) error {
	//Busco los datos del producto en cache
	client := configs.Client()
	result, err := client.Get(id).Result()
	if err != nil {
		fmt.Println("no se encontro el registro en caché, buscando en el microservicio correspondiente...")
	} else {
		product, err := json.Marshal(result)
		if err != nil {
			return custerror.NewCustom(500, "no se pudo converir el datos del producto en cache")
		}
		fmt.Println("retornando de cache:", product)
		uncuotedProduct, err := strconv.Unquote(string(product))
		if err != nil {
			fmt.Println(err)
		}
		return json.Unmarshal([]byte(uncuotedProduct), target)
	}
	req, err := http.NewRequest("GET", env.Get().CatalogURL+"/v1/articles/"+id, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "bearer "+token)
	response, err := http.DefaultClient.Do(req)
	if err != nil || response.StatusCode != 200 {
		return err
	}
	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&target)
	if err != nil {
		return err
	}
	cacheProd, err := json.Marshal(&target)
	if err != nil {
		return err
	}
	//guardamos el producto en cache una vez encontrado
	fmt.Println("json a cachear:", cacheProd)

	err = client.Set(id, cacheProd, 1*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}
