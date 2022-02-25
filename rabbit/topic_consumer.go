package rabbit

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"statsv0/controllers"
	"statsv0/models"
	"statsv0/services"
	"statsv0/tools/env"

	"github.com/streadway/amqp"
)

/**
 * @api {topic} stats/orders Escucha las ordenes nuevas
 * @apiGroup RabbitMQ GET
 *
 * @apiDescription Escucha de mensajes order_placed desde orders.
 *
 * @apiSuccessExample {json} Mensaje
 *    {
 *"type": "order-placed",
 *"message" : {
 *    "cartId": "{cartId}",
 *    "orderId": "{orderId}"
 *    "articles": [{
 *         "articleId": "{article id}"
 *         "quantity" : {quantity}
 *     }, ...]
 *   }
 *}
 */

// Init se queda escuchando topics de ordenes
func InitOrders() {
	go func() {
		fmt.Println("RabbitMQ escuchando ordenes Y logouts")
		for {
			listenOrders()
			time.Sleep(5 * time.Second)
		}
	}()
}

func listenOrders() error {
	conn, err := amqp.Dial(env.Get().RabbitURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	chn, err := conn.Channel()
	if err != nil {
		return err
	}
	defer chn.Close()

	err = chn.ExchangeDeclare(
		"sell_flow", // name
		"topic",     // type
		false,       // durable
		false,       // auto-deleted
		false,       // internal
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		println("no se pudo declarar el exchange")
		print(err)
		return err
	}

	queue, err := chn.QueueDeclare(
		"order", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		println("no se pudo declarar la cola order")
		return err
	}

	err = chn.QueueBind(
		queue.Name,     // queue name
		"order_placed", // routing key
		"sell_flow",    // exchange
		false,
		nil)
	if err != nil {
		println("no se pudo bindear las colas")
		return err
	}

	mgs, err := chn.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return err
	}

	fmt.Println("RabbitMQ conectado")

	// Por cada orden recibida, tomamos el id de la orden,  y los ids de los productos pertenecientes a estas.
	// Luego, se guarda la orden recibida en el esquema "Greatest Orders", y cada uno de articulos de estas en "Greatest Products",
	// asi como tambien, por cada articulo vendido y segun el mes año y dia de semana, se guarda una instancia de "SellsPerDay".
	// Finalmente, se incrementa en 1 dependiendo de la hora en la que fue enviada la orden, la hora pico con el esquema "Peak Hour".
	// De esta maera se puebla la base de datos con todos los datos necesarios para obtener las estadisticas cuando estas sean necesarias.
	go func() {
		for d := range mgs {
			log.Output(1, "Orden recibida")
			println("Orden recibida")

			newMessage := &models.MsgOrder{}
			err = json.Unmarshal(d.Body, newMessage)
			id := newMessage.Message.OrderId
			year, month, _ := time.Now().Date()
			weekday := int(time.Now().Weekday())
			articles := newMessage.Message.Articles
			articlesQ := GetTotalArticlesQuantity(newMessage.Message.Articles)

			if err == nil {
				fmt.Println(newMessage.Type)
				fmt.Println(newMessage.Message)
				if newMessage.Type == "order-placed" {
					controllers.GuardarGreatestOrder(id, int(month), year, articlesQ) // Guardamos el id de la orden con su cantidad de articulos
					controllers.IncPeakHour(year, int(month), time.Now().Hour())      //Se incrementan (o se crea una nueva instancia) las ventas segun la hora
					for _, article := range articles {
						controllers.GuardarGreatestProduct(article, int(month), year)   //Se guarda el articulo con sus datos (incluida la cantidad vendida)
						services.GuardarSellsPerDay(article, int(month), year, weekday) //Se guarda el articulo en base al dia de la semana
					}
				}
			}
		}
	}()
	fmt.Print("Closed connection: ", <-conn.NotifyClose(make(chan *amqp.Error)))

	return nil
}

//Obtiene la cantidad total de articulos a partir de la lista de articulos presente en cada orden en caso de que esta posea mas de 1 articulo en ella.
func GetTotalArticlesQuantity(arts []models.Article) int {
	sum := 0
	for _, art := range arts {
		sum = sum + art.Quantity
	}
	print("cantidad:", sum)
	return sum
}
