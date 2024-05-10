package models

var Category = [4]string{
	"Clothing",
	"Accessories",
	"Footwear",
	"Beverages",
}

type Product struct {
	Id string `db:"id" json:"id"`
	// UserId      string `json:"userId" db:"userId"`
	Name        string  `db:"name" json:"name"`
	Sku         string  `db:"sku" json:"sku"`
	Category    string  `db:"category" json:"category"`
	ImageUrl    string  `db:"imageUrl" json:"imageUrl"`
	Notes       string  `db:"notes" json:"notes"`
	Price       int     `db:"price" json:"price"`
	Stock       int     `db:"stock" json:"stock"`
	Location    string  `db:"location" json:"location"`
	IsAvailable bool    `db:"isAvailable" json:"isAvailable"`
	CreatedAt   string  `db:"createdAt" json:"createdAt"`
	UpdatedAt   string  `db:"updatedAt" json:"updatedAt"`
	DeletedAt   *string `db:"deletedAt" json:"deletedAt,omitempty"`
}

type ProductCheckout struct {
	Product
	Quantity int `json:"quantity"`
}