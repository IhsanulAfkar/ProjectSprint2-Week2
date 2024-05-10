package db

import (
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)
var db *sqlx.DB
var err error
func Init() {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PARAMS"),
	)
	db, err = sqlx.Connect("pgx", connStr)
	if err != nil {
		panic(err.Error())
	}
	db.DB.SetMaxOpenConns(100)
	db.DB.SetMaxIdleConns(10)
	// Migrations
	// migrate -database "postgresql://root:root@localhost:5432/eniqilo?sslmode=disable" -path db/migrations up
}
func CreateConn() *sqlx.DB{
	return db
}