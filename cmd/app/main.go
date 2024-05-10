package main

import (
	"encoding/json"
	"fmt"
	"http/standarlibary/models"
	"log"
	"net/http"
	"regexp"

	"github.com/gosimple/slug"
)

var (
	CarRgx = regexp.MustCompile(`^/cars/*$`)
	CarRgxID = regexp.MustCompile(`^/cars/([a-zA-Z0-9]+(?:-[0-9]+))$`)
)

func main() {
	store := models.NewRedisHandler()
	carsHandler := newCarsHandler(store)
	if store == nil{
		log.Fatalf("ERROR:\nCan't connect to Redis\n")
	}
	mux := http.NewServeMux()

	mux.Handle("/", &homeHandler{})
	mux.Handle("/cars", carsHandler)
	mux.Handle("/cars/", carsHandler)
	http.ListenAndServe(":8080", mux)
}

type carStore interface {
	Add(name string, car models.Car) error
	Get(name string) (models.Car, error)
	Update(name string, car models.Car) error
	List() (map[string]models.Car, error)
	Remove(name string) error
}

type CarsHandler struct{
	store carStore
}

func newCarsHandler (c carStore) *CarsHandler {
	return &CarsHandler {
		store: c,
	}
}

func (c *CarsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	switch{
		case r.Method == http.MethodPost && CarRgx.MatchString(r.URL.Path):
			c.CreateCar(w, r)
			return
		case r.Method == http.MethodGet && CarRgx.MatchString(r.URL.Path):
			c.ListCar(w, r)
			return
		case r.Method == http.MethodGet && CarRgxID.MatchString(r.URL.Path):
			c.GetCar(w, r)
			return
		case (r.Method == http.MethodPut || r.Method == http.MethodPatch)&& CarRgxID.MatchString(r.URL.Path):
			c.UpdateCar(w,r)
			return
		case r.Method == http.MethodDelete && CarRgxID.MatchString(r.URL.Path):
			c.DeleteCar(w,r)
			return
		default:
			NotFoundHandler(w,r)
			return
	}
}

type homeHandler struct{}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	fmt.Println("Home")
	w.Write([]byte("This is my home page"))
}

func (h *CarsHandler) CreateCar(w http.ResponseWriter, r *http.Request){
	var car models.Car

	if err := json.NewDecoder(r.Body).Decode(&car); err!=nil{
		fmt.Println(car)
		InternalServerErrorHandler(w,r)
		return
	}

	if car.Model == "" {
        BadRequestHandler(w,r)
		return
    }

    if car.Year < 1000 || car.Year > 9999 {
        BadRequestHandler(w,r)
		return
    }


	ID := car.Model+"-"+fmt.Sprint(car.Year)
	fmt.Printf("ID: %v\n",ID)
	resourceID := slug.Make(ID)
	if err := h.store.Add(resourceID,car); err != nil{
		InternalServerErrorHandler(w,r)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"key": "`+resourceID+`"}`))
}

func (h *CarsHandler) ListCar(w http.ResponseWriter, r *http.Request){
	resources, err := h.store.List()
	if err != nil {
		InternalServerErrorHandler(w,r)
		return
	}
	
	json, err := json.Marshal(resources)
	if err != nil {
		InternalServerErrorHandler(w,r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func (h *CarsHandler) GetCar(w http.ResponseWriter, r *http.Request){
	matches := CarRgxID.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		InternalServerErrorHandler(w,r)
		return
	}

	car, err := h.store.Get(matches[1])
	if err != nil {
		if err == models.ErrNotFound{
			NotFoundHandler(w,r)
			return
		}
		InternalServerErrorHandler(w,r)
		return
	}

	json, err := json.Marshal(car)
	if err != nil {
		InternalServerErrorHandler(w,r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(json)
}
func (h *CarsHandler) UpdateCar(w http.ResponseWriter, r *http.Request){
	matches := CarRgxID.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		InternalServerErrorHandler(w,r)
		return
	}

	var car models.Car
	if err := json.NewDecoder(r.Body).Decode(&car); err != nil{
		InternalServerErrorHandler(w,r)
		return
	}
	
	if err := h.store.Update(matches[1],car); err != nil {
		if err == models.ErrNotFound{
			NotFoundHandler(w,r)
			return
		}
		InternalServerErrorHandler(w,r)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"key": "`+matches[1]+`"}`))
}
func (h *CarsHandler) DeleteCar(w http.ResponseWriter, r *http.Request){
	matches := CarRgxID.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		InternalServerErrorHandler(w,r)
		return
	}

	if err := h.store.Remove(matches[1]); err != nil {
		if err == models.ErrNotFound{
			NotFoundHandler(w,r)
			return
		}
		InternalServerErrorHandler(w,r)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
}
func BadRequestHandler(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("400 Bad Request"))
}


func NotFoundHandler(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}
