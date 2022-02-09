package models

type SellsPerDay struct {
	ProductId string `json:"productid,omitempty"`
	Month     int    `json:"month,omitempty" validate:"required"`
	Year      int    `json:"year,omitempty" validate:"required"`
	Weekday   int    `json:"weekday,omitempty" validate:"required"`
	Quantity  int    `json:"quantity,omitempty" validate:"required"`
}

type BestProductDay struct {
	Day         int     `json:"day,omitempty"`
	Total       int     `json:"total,omitempty"`
	Id          string  `json:"id,omitempty"`
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Image       string  `json:"image,omitempty"`
	Price       float64 `json:"price,omitempty"`
	Stock       int     `json:"stock,omitempty"`
	Updated     string  `json:"updated,omitempty"`
	Created     string  `json:"created,omitempty"`
	Enabled     string  `json:"enabled,omitempty"`
}
type AggSellsPerDay struct {
	ID       Ide `bson:"_id,omitempty"`
	Quantity int `bson:"articlequantity,omitempty"`
}
type Ide struct {
	Weekday int `json:"weekday,omitempty" `
	Month   int `json:"month,omitempty" `
	Year    int `json:"year,omitempty" `
}
type Product struct {
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Image       string  `json:"image,omitempty"`
	Price       float64 `json:"price,omitempty"`
	Stock       float64 `json:"stock,omitempty"`
	Enabled     bool    `json:"enabled,omitempty"`
}
type SellsPerDayProduct struct {
	Weekday     int     `json:"weekday,omitempty" `
	Month       int     `json:"month,omitempty" `
	Year        int     `json:"year,omitempty" `
	Quantity    int     `bson:"articlequantity,omitempty"`
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Image       string  `json:"image,omitempty"`
	Price       float64 `json:"price,omitempty"`
	Stock       float64 `json:"stock,omitempty"`
	Enabled     bool    `json:"enabled,omitempty"`
}
