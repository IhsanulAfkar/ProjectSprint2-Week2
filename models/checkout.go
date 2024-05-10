package models

import "Week2/forms"

type Checkout struct {
	Id         string `db:"id" json:"id"`
	CustomerId string    `json:"customerId" db:"customerId"`
	Total      int    `json:"total" db:"total"`
	Paid int  `json:"paid" db:"paid"`
	Change         int             `json:"change" db:"change"`
	CreatedAt  string `db:"createdAt" json:"createdAt"`
}

type CheckoutItem struct {
	Id         string `db:"id" json:"id"`
	CheckoutId string    `json:"checkoutId" db:"checkoutId"`
	ProductId  string `json:"productId" db:"productId"`
	Quantity   int    `json:"quantity" db:"quantity"`
	Price      int    `json:"price" db:"price"`
	CreatedAt  string `db:"createdAt" json:"createdAt"`
}

type GetCheckout struct {
	TransactionId string `json:"transactionId"`
	CustomerId string    `json:"customerId"`
	ProductDetails []forms.ProductDetail `json:"productDetails"`
	Paid           int             `json:"paid"`
	Change         int             `json:"change"`
}