package models

type Staff struct {
	Id        string  `db:"id" json:"id"`
	Phone     string  `json:"phoneNumber" db:"phone"`
	Name      string  `db:"name" json:"name"`
	Password  string  `db:"password" json:"password"`
	CreatedAt string  `db:"createdAt" json:"createdAt"`
	UpdatedAt string  `db:"updatedAt" json:"updatedAt"`
	DeletedAt *string `db:"deletedAt" json:"deletedAt"`
}