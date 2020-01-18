package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type OperationType string

const (
	Create OperationType = "CREATE"
	Read   OperationType = "READ"
	Update OperationType = "UPDATE"
	Delete OperationType = "DELETE"
)

type DBError struct {
	Operation OperationType
	Err       error
}

func (e *DBError) Error() string {
	return fmt.Sprintf("[DBError][%s]errors:%w", e.Operation, e.Err)
}

type CarInventory struct {
	Vin    string `json:"vin"`
	Model  string `json:"model"`
	Make   string `json:"make"`
	Year   int    `json:"year"`
	MSRP   int    `json:"msrp"`
	Status string `json:"status"`
	Booked string `json:"booked"`
	Listed string `json:"listed"`
}

type Data interface{}

type ResponseError struct {
	Code    int    `json:code`
	Message string `json:message`
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("[ResponseError](%d)[%s]", e.Code, e.Message)
}

type APIResponse struct {
	Success bool          `json:"success"`
	Data    Data          `json:"data"`
	Err     ResponseError `json:"err"`
}

const DATA_FILENAME = "data.json"
const SERVER_PORT = 8888

var carInventories []CarInventory

func createCarInventory(filename string, data CarInventory) (CarInventory, error) {
	carInventories, readDBError := readCarInventories(DATA_FILENAME)
	if readDBError != nil {
		return data, &DBError{
			Operation: Create,
			Err:       readDBError,
		}
	}

	carInventories = append(carInventories, data)
	byteData, jsonError := json.MarshalIndent(carInventories, " ", "")
	if jsonError != nil {
		return data, &DBError{
			Operation: Create,
			Err:       jsonError,
		}
	}

	fileError := ioutil.WriteFile(filename, byteData, 0644)
	if fileError != nil {
		dbError := &DBError{
			Operation: Create,
			Err:       fileError,
		}
		return data, dbError
	}

	return data, nil
}

func readCarInventories(filename string) ([]CarInventory, error) {
	byteCarInventories, fileError := ioutil.ReadFile(filename)
	if fileError != nil {
		return nil, &DBError{
			Operation: Read,
			Err:       fileError,
		}
	}

	jsonError := json.Unmarshal(byteCarInventories, &carInventories)
	if jsonError != nil {
		return nil, &DBError{
			Operation: Read,
			Err:       jsonError,
		}
	}

	return carInventories, nil
}

func deleteCarInventory(filename string, vin string) (string, error) {
	carInventories, readDBError := readCarInventories(DATA_FILENAME)
	if readDBError != nil {
		return vin, &DBError{
			Operation: Delete,
			Err:       readDBError,
		}
	}

	for i, carInventory := range carInventories {
		if carInventory.Vin == vin {
			carInventories = append(carInventories[:i], carInventories[i+1:]...)
			break
		}
	}

	byteData, jsonError := json.MarshalIndent(carInventories, " ", "")
	if jsonError != nil {
		return vin, &DBError{
			Operation: Delete,
			Err:       jsonError,
		}
	}

	fileError := ioutil.WriteFile(filename, byteData, 0644)
	if fileError != nil {
		return vin, &DBError{
			Operation: Delete,
			Err:       fileError,
		}
	}

	return vin, nil
}

func loadCarInventoryData() {
	data, err := ioutil.ReadFile("data.json")
	if err != nil {
		log.Fatal("car inventories are not initialized")
	}

	json.Unmarshal(data, &carInventories)
}

func handleReadInventoryCars(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	carInventories, err := readCarInventories(DATA_FILENAME)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		carInventories = []CarInventory{}
	}

	json.NewEncoder(w).Encode(&APIResponse{
		Success: true,
		Data:    carInventories,
	})
}

func handleCreateInventoryCar(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var c CarInventory
	json.NewDecoder(r.Body).Decode(&c)
	if c.Vin == "" {
		log.Fatal("Vin Field Required!")
	}

	_, dbError := createCarInventory(DATA_FILENAME, c)
	if dbError != nil {
		json.NewEncoder(w).Encode(&APIResponse{
			Success: false,
			Err: ResponseError{
				Code:    20001,
				Message: dbError.Error(),
			},
		})
	}

	json.NewEncoder(w).Encode(&APIResponse{
		Success: true,
		Data:    c,
	})
}

func handleDeleteInventoryCar(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	vin := ps.ByName("vin")

	_, dbError := deleteCarInventory(DATA_FILENAME, vin)
	if dbError != nil {
		w.WriteHeader(500)
	}

	json.NewEncoder(w).Encode(&APIResponse{
		Success: true,
		Data:    vin,
	})
}

func main() {
	router := httprouter.New()

	router.GET("/inventory/cars", handleReadInventoryCars)
	router.POST("/inventory/cars", handleCreateInventoryCar)
	router.DELETE("/inventory/cars/:vin", handleDeleteInventoryCar)

	fmt.Printf("api server listening on %d\n", SERVER_PORT)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", SERVER_PORT), &Server{router}))
}

type Server struct {
	r *httprouter.Router
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	header.Set("Content-Type", "application/json")
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	s.r.ServeHTTP(w, r)
}
