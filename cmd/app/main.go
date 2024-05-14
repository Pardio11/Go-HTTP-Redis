package main

import (
	"http/standarlibary/handlers"
	"http/standarlibary/models"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
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
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000/"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowedHeaders: []string{"X-Requested-With", "Content-Type", "Authorization"},
	})
	http.ListenAndServe(go_port, c.Handler(mux))
}





