package models

type PeakHour struct {
	Hour          int `json:"hour,omitempty"`
	Month         int `json:"month,omitempty" validate:"required"`
	Year          int `json:"year,omitempty" validate:"required"`
	OrderQuantity int `json:"orderQuantity,omitempty" validate:"required"`
}
