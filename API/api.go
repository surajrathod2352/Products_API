package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"apiwdb/DB"
	"apiwdb/products"
	"apiwdb/receipt"
)

const basePath = "/api"

func main() {
	database.SetupDatabase()
	receipt.SetupRoutes(basePath)
	product.SetupRoutes(basePath)
	log.Fatal(http.ListenAndServe(":5000", nil))
}
