package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"http/standarlibary/handlers"
	"http/standarlibary/models"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gosimple/slug"
	"github.com/stretchr/testify/assert"
)

func readTestData(t *testing.T, name string) []byte {
	t.Helper()
	content, err := os.ReadFile("./"+name)
	if err != nil {
		t.Errorf("Couldn't read: %v",name)
	}
	return content
}

func TestCarsHandler(t *testing.T){

	store := models.NewRedisHandler()
	carsHandler := handlers.NewCarsHandler(store)

	carJson := readTestData(t, "cars.json")
	carReader := bytes.NewReader(carJson)
	var car models.Car
	if err := json.NewDecoder(carReader).Decode(&car); err != nil{
		fmt.Println("Error decoding JSON:", err)
		return
	}

	carReader = bytes.NewReader(carJson)
	id:=slug.Make(car.Model+"-"+fmt.Sprint(car.Year))

	// RESET
    req := httptest.NewRequest(http.MethodDelete, "/cars/"+id, nil)
    w := httptest.NewRecorder()
    carsHandler.ServeHTTP(w, req)

    res := w.Result()
    defer res.Body.Close()
    if res.StatusCode != 200 && res.StatusCode != 404 {
		t.Errorf("Unexpected status code: %d", res.StatusCode)
	}
	// RESET

	//CREATE//
	saved, _ := store.List()
    size := len(saved)
	expectedReturn := `{"key": "`+id+`"}`
	req = httptest.NewRequest(http.MethodPost, "/cars", carReader)
	w = httptest.NewRecorder()
	carsHandler.ServeHTTP(w,req)
	
	res = w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
	saved, _ = store.List()
    assert.Len(t, saved, size+1)
	assert.JSONEq(t,expectedReturn,string(data))
	assert.Equal(t, 200, res.StatusCode)
	//CREATE//
	
	// GET //
    req = httptest.NewRequest(http.MethodGet, "/cars/"+id, nil)
    w = httptest.NewRecorder()
    carsHandler.ServeHTTP(w, req)

    res = w.Result()
    defer res.Body.Close()
    assert.Equal(t, 200, res.StatusCode)

    data, err = io.ReadAll(res.Body)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
    assert.JSONEq(t, string(carJson), string(data))
	// GET //


	// LIST //
    req = httptest.NewRequest(http.MethodGet, "/cars", nil)
    w = httptest.NewRecorder()
    carsHandler.ServeHTTP(w, req)

    res = w.Result()
    defer res.Body.Close()
    assert.Equal(t, 200, res.StatusCode)
	// LIST //

	// UPDATE //
	carReader = bytes.NewReader(carJson)
	req = httptest.NewRequest(http.MethodPut, "/cars/A-2000", carReader)
	w = httptest.NewRecorder()
	carsHandler.ServeHTTP(w,req)
	
	res = w.Result()
	defer res.Body.Close()
	assert.Equal(t, 404, res.StatusCode)
	// UPDATE //

	// UPDATE //
	carJson = readTestData(t, "carsUpdate.json")
	carReader = bytes.NewReader(carJson)
	req = httptest.NewRequest(http.MethodPut, "/cars/"+id, carReader)
	w = httptest.NewRecorder()
	carsHandler.ServeHTTP(w,req)
	
	res = w.Result()
	defer res.Body.Close()
	data, err = io.ReadAll(res.Body)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
	assert.JSONEq(t,expectedReturn,string(data))
	assert.Equal(t, 200, res.StatusCode)
	// UPDATE //

	// DELETE //
    req = httptest.NewRequest(http.MethodDelete, "/cars/"+id, nil)
    w = httptest.NewRecorder()
    carsHandler.ServeHTTP(w, req)

    res = w.Result()
    defer res.Body.Close()
    assert.Equal(t, 200, res.StatusCode)
	// DELETE //


}