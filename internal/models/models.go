package models

import "time"

type Cryptocurrency struct {
	Coin      string    `json:"coin" db:"coin"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type PriceHistory struct {
	ID        int       `json:"id" db:"id"`
	Coin      string    `json:"coin" db:"coin"`
	Price     float64   `json:"price" db:"price"`
	Timestamp int64     `json:"timestamp" db:"timestamp"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type PriceRequest struct {
	Coin      string `json:"coin" binding:"required"`
	Timestamp int64  `json:"timestamp" binding:"required"`
}

type AddCoinRequest struct {
	Coin string `json:"coin" binding:"required"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
