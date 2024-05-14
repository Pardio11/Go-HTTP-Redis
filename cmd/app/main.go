package main

import (
	"http/standarlibary/handlers"
	"http/standarlibary/models"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)



func main() {
	store := models.NewRedisHandler()
	carsHandler := handlers.NewCarsHandler(store)
	if store == nil{
		log.Fatalf("ERROR:\nCan't connect to Redis\n")
	}
	mux := http.NewServeMux()

	mux.Handle("/", &handlers.HomeHandler{})
	mux.Handle("/cars", carsHandler)
	mux.Handle("/cars/", carsHandler)
	godotenv.Load()
	go_port := os.Getenv("GO_PORT")
	http.ListenAndServe(go_port, mux)
}





