package models

type GreatestOrders struct {
	Id              string `json:"id,omitempty"`
	Month           int    `json:"month,omitempty" validate:"required"`
	Year            int    `json:"year,omitempty" validate:"required"`
	ArticleQuantity int    `json:"articleQuantity,omitempty" validate:"required"`
}
type GetGO struct {
	Id             string  `json:"id,omitempty"`
	Status         string  `json:"status,omitempty"`
	Total          float64 `json:"total,omitempty"`
	Created        string  `json:"created,omitempty"`
	TotalProductos int     `json:"totalProductos,omitempty"`
}
