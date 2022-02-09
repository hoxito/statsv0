package models

type MsgOrder struct {
	Type    string `json:"type"`
	Message Order  `json:"message"`
}
type Order struct {
	CartId   string    `json:"cartId"`
	OrderId  string    `json:"orderId"`
	Articles []Article `json:"articles"`
}
type Article struct {
	ArticleId string `json:"articleId"`
	Quantity  int    `json:"quantity"`
}

type Orden struct {
	Id           string     `json:"id"`
	Status       string     `json:"status"`
	CartId       string     `json:"cartId"`
	TotalPrice   float64    `json:"totalPrice"`
	TotalPayment float64    `json:"totalPayment"`
	Updated      string     `json:"updated"`
	Created      string     `json:"created"`
	Articulos    []Articulo `json:"articles"`
	Payment      []float64  `json:"payment"`
}
type Articulo struct {
	ArticleId    string  `json:"articleId"`
	Quantity     int     `json:"quantity"`
	UnitaryPrice float64 `json:"unitaryPrice"`
	Validated    bool    `json:"validated"`
	Valid        bool    `json:"valid"`
}
