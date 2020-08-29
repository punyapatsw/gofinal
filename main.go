package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/punyapatsw/gofinal/database"
	"github.com/punyapatsw/gofinal/middleware"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB
var err error

func init() {
	url := os.Getenv("DATABASE_URL")
	db, err = sql.Open("postgres", url)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("Connected")
}

func setupRouter() *gin.Engine {
	h := database.Handler{DB: db}
	h.CreateTableHandler()
	r := gin.Default()
	r.Use(middleware.Auth)
	r.POST("/customers", h.CreateCustomerHandler)
	r.GET("/customers/:id", h.GetCustomerHandler)
	r.GET("/customers/", h.GetAllCustomerHandler)
	r.PUT("/customers/:id", h.UpdateCustomerHandler)
	r.DELETE("/customers/:id", h.DeleteCustomerHandler)

	return r
}

func main() {
	fmt.Println("customer service")
	//run port ":2009"
	r := setupRouter()
	r.Run(":2009")
	defer db.Close()
}
