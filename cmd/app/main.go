package main

import (
	"http/standarlibary/handlers"
	"http/standarlibary/models"
	"log"
	"net/http"
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
	http.ListenAndServe(":8080", mux)
}





