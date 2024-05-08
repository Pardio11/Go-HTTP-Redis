package main

import (
	"net/http"
	"regexp"
)

var (
	CarRgx = regexp.MustCompile(`^/recipes/*$`)
	CarRgxID = regexp.MustCompile(`^/recipes/([a-z0-9]+(?:-[a-z0-9]+)+)$`)
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/", &homeHandler{})
	mux.Handle("/cars", &CarsHandler{})
	mux.Handle("/cars/", &CarsHandler{})
	http.ListenAndServe(":8080", mux)
}

type CarsHandler struct{}

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
			return
	}
}


type homeHandler struct{}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("This is my home page"))
}

func (c *CarsHandler) CreateCar(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("This is my create car Endpoint"))
}
func (c *CarsHandler) ListCar(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("This is my list car Endpoint"))
}
func (c *CarsHandler) GetCar(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("This is my get car Endpoint"))
}
func (c *CarsHandler) UpdateCar(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("This is my upadate car Endpoint"))
}
func (c *CarsHandler) DeleteCar(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("This is my delete car Endpoint"))
}
