package main

import (
	"Week2/db"
	"Week2/routes"

	_ "github.com/joho/godotenv/autoload"
)
func main(){
	db.Init()
	r := routes.Init()
	r.Run(":8080")
}