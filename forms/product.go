package forms

type ProductCreate struct {
	Name        string `json:"name"`
	Sku         string `json:"sku"`
	Category    string `json:"category"`
	ImageUrl    string `json:"imageUrl"`
	Notes       string `json:"notes"`
	Price       int    `json:"price"`
	Stock       int    `json:"stock"`
	Location    string `json:"location"`
	IsAvailable bool   `json:"isAvailable"`
}
type ProductDetail struct {
	ProductId string `json:"productId"`
	Quantity  int    `json:"quantity"`
}
type Checkout struct {
	CustomerId     string          `json:"customerId"`
	ProductDetails []ProductDetail `json:"productDetails"`
	Paid           int             `json:"paid"`
	Change         int             `json:"change"`
}