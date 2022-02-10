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
 * @apiDescription Escucha de mensajes logout desde orders.
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
					controllers.GuardarGreatestOrder(id, int(month), year, articlesQ)
					controllers.IncPeakHour(year, int(month), time.Now().Hour())
					for _, article := range articles {
						controllers.GuardarGreatestProduct(article, int(month), year)
						services.GuardarSellsPerDay(article, int(month), year, weekday)
					}
				}
			}
		}
	}()
	fmt.Print("Closed connection: ", <-conn.NotifyClose(make(chan *amqp.Error)))

	return nil
}

func GetTotalArticlesQuantity(arts []models.Article) int {
	sum := 0
	for _, art := range arts {
		sum = sum + art.Quantity
	}
	print("cantidad:", sum)
	return sum
}
