package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// Order Полная модель заказа
// @Descriptionn Модель описывает данные возвращаемого заказа
type Order struct {
	OrderUID          uuid.UUID `json:"order_uid" validate:"required"`
	TrackNumber       string    `json:"track_number" validate:"required"`
	Entry             string    `json:"entry" validate:"required"`
	Delivery          Delivery  `json:"delivery" validate:"required"`
	Payment           Payment   `json:"payment" validate:"required"`
	Items             []Item    `json:"items" validate:"required"`
	Locale            string    `json:"locale" validate:"required"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id" validate:"required"`
	DeliveryService   string    `json:"delivery_service" validate:"required"`
	Shardkey          string    `json:"shardkey" validate:"required"`
	SmID              int       `json:"sm_id" validate:"gte=0"`
	DateCreated       time.Time `json:"date_created" validate:"required"`
	OofShard          string    `json:"oof_shard"`
}

// Delivery Модель доставки
// @Description Модель описывает информацию о доставщике
type Delivery struct {
	ID      int    `json:"-"`
	Name    string `json:"name" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
	Zip     string `json:"zip"`
	City    string `json:"city" validate:"required"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

// Payment
// @Description Модель описывает информацию о платеже
type Payment struct {
	ID           int       `json:"-"`
	Transaction  uuid.UUID `json:"transaction" validate:"required"`
	RequestID    string    `json:"request_id"`
	Currency     string    `json:"currency" validate:"required"`
	Provider     string    `json:"provider" validate:"required"`
	Amount       int       `json:"amount" validate:"gte=0"`
	PaymentDt    int64     `json:"payment_dt" validate:"gt=0"`
	Bank         string    `json:"bank" validate:"required"`
	DeliveryCost int       `json:"delivery_cost" validate:"gte=0"`
	GoodsTotal   int       `json:"goods_total" validate:"gte=0"`
	CustomFee    int       `json:"custom_fee" validate:"gte=0"`
}

// Item
// @Description Модель описывает информацию о товаре в заказе
type Item struct {
	ChrtID      int    `json:"chrt_id" validate:"required"`
	TrackNumber string `json:"track_number" validate:"required"`
	Price       int    `json:"price" validate:"gte=0"`
	Rid         string `json:"rid" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Sale        int    `json:"sale" validate:"gte=0,lte=100"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price" validate:"gte=0"`
	NmID        int    `json:"nm_id" validate:"gte=0"`
	Brand       string `json:"brand"`
	Status      int    `json:"status" validate:"gt=0"`
}
