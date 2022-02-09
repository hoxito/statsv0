package models

type GreatestProducts struct {
	ProductId string `json:"productId,omitempty"`
	Month     int    `json:"month,omitempty" validate:"required"`
	Year      int    `json:"year,omitempty" validate:"required"`
	Sells     int    `json:"articleQuantity,omitempty" validate:"required"`
}
type AggGreatestProducts struct {
	ID    GPId `bson:"_id,omitempty"`
	Sells int  `bson:"ventasTotales,omitempty"`
}
type GPId struct {
	ProductId string `json:"productid,omitempty" `
	Month     int    `json:"month,omitempty" `
	Year      int    `json:"year,omitempty" `
}
type GPProduct struct {
	ProductId   string  `json:"productid,omitempty" `
	Month       int     `json:"month,omitempty" `
	Year        int     `json:"year,omitempty" `
	TotalSells  int     `bson:"articlequantity,omitempty"`
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Image       string  `json:"image,omitempty"`
	Price       float64 `json:"price,omitempty"`
	Stock       float64 `json:"stock,omitempty"`
	Enabled     bool    `json:"enabled,omitempty"`
}
