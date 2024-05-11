package models

type Customer struct {
	// PkId        int    `db:"pkId" json:"pkId"`
	Id        string  `db:"id" json:"userId"`
	Phone     string  `db:"phone" json:"phoneNumber"`
	Name      string  `db:"name" json:"name"`
	CreatedAt string  `db:"createdAt" json:"createdAt"`
	UpdatedAt string  `db:"updatedAt" json:"updatedAt"`
	DeletedAt *string `db:"deletedAt" json:"deletedAt"`
}