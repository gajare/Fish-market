package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gajare/Fish-market/controller"
	"github.com/gajare/Fish-market/db"
	"github.com/gajare/Fish-market/router"
	"github.com/gajare/Fish-market/service"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	db.Connect()

	svc := service.NewUserService()
	uc := controller.NewUserController(svc)
	r := router.New(uc)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("user service listening on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
