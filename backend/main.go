package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/go-playground/validator"

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
	Vin    string `json:"vin" validate:"required"`
	Model  string `json:"model" validate:"required"`
	Make   string `json:"make" validate:"required"`
	Year   int    `json:"year" validate:"required"`
	MSRP   int    `json:"msrp" validate:"required"`
	Status string `json:"status" validate:"required"`
	Booked string `json:"booked" validate:"required"`
	Listed string `json:"listed" validate:"required"`
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

var (
	readFromCarInventoryDB = readFromFile(DATA_FILENAME)
	commitToCarInventoryDB = commitToFile(DATA_FILENAME, 0644)
	createCarInventoryDB   = func() error {
		initCarInventory, jsonError := encodeJSON([]CarInventory{})
		if jsonError != nil {
			return jsonError
		}
		return commitToCarInventoryDB(initCarInventory)
	}
)

func readFromFile(filename string) func() ([]byte, error) {
	return func() ([]byte, error) {
		return ioutil.ReadFile(filename)
	}
}

func commitToFile(filename string, perm os.FileMode) func([]byte) error {
	return func(data []byte) error {
		return ioutil.WriteFile(filename, data, perm)
	}
}

func encodeJSON(data interface{}) ([]byte, error) {
	return json.MarshalIndent(data, " ", "")
}

func decodeJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func createCarInventory(data CarInventory) (CarInventory, error) {
	carInventories, readDBError := readCarInventories()
	if readDBError != nil {
		return data, &DBError{
			Operation: Create,
			Err:       readDBError,
		}
	}

	carInventories = append(carInventories, data)
	byteData, jsonError := encodeJSON(carInventories)
	if jsonError != nil {
		return data, &DBError{
			Operation: Create,
			Err:       jsonError,
		}
	}

	commitError := commitToCarInventoryDB(byteData)
	if commitError != nil {
		return data, &DBError{
			Operation: Create,
			Err:       commitError,
		}
	}

	return data, nil
}

func readCarInventories() ([]CarInventory, error) {
	var carInventories []CarInventory
	data, notExistCarInventoryDB := readFromCarInventoryDB()
	if notExistCarInventoryDB != nil {
		createDBError := createCarInventoryDB()
		if createDBError != nil {
			return nil, &DBError{
				Operation: Read,
				Err:       createDBError,
			}
		}
		data, _ = readFromCarInventoryDB()
	}

	jsonError := decodeJSON(data, &carInventories)
	if jsonError != nil {
		return nil, jsonError
	}

	return carInventories, nil
}

func deleteCarInventory(vin string) (string, error) {
	carInventories, readError := readCarInventories()
	if readError != nil {
		return vin, &DBError{
			Operation: Delete,
			Err:       readError,
		}
	}

	for i, carInventory := range carInventories {
		if carInventory.Vin == vin {
			carInventories = append(carInventories[:i], carInventories[i+1:]...)
			break
		}
	}

	byteData, jsonError := encodeJSON(carInventories)
	if jsonError != nil {
		return vin, &DBError{
			Operation: Delete,
			Err:       jsonError,
		}
	}

	fileError := commitToCarInventoryDB(byteData)
	if fileError != nil {
		return vin, &DBError{
			Operation: Delete,
			Err:       fileError,
		}
	}

	return vin, nil
}

func parseHTTPBody(r *http.Request, d interface{}) {
	json.NewDecoder(r.Body).Decode(&d)
}

func validateCarInventorySchema(c CarInventory) bool {
	v := validator.New()
	err := v.Struct(c)
	if err != nil {
		return false
	}
	return true
}

func responseJSON(w http.ResponseWriter, r *http.Request, apiResponse APIResponse) {
	json.NewEncoder(w).Encode(apiResponse)
}

func responseInvalidParameter(w http.ResponseWriter, r *http.Request) {
	apiResponse := APIResponse{
		Success: false,
		Err: ResponseError{
			Code:    30001,
			Message: "Some parameters are invalid.",
		},
	}
	w.WriteHeader(http.StatusBadRequest)
	responseJSON(w, r, apiResponse)
}

func responseDBError(w http.ResponseWriter, r *http.Request) {
	apiResponse := APIResponse{
		Success: false,
		Err: ResponseError{
			Code:    20001,
			Message: "Internal Server Error",
		},
	}
	w.WriteHeader(http.StatusInternalServerError)
	responseJSON(w, r, apiResponse)
}

func handleReadInventoryCars(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	carInventories, err := readCarInventories()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		carInventories = []CarInventory{}
	}

	responseJSON(w, r, APIResponse{
		Success: true,
		Data:    carInventories,
	})
}

func handleCreateInventoryCar(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var c CarInventory
	parseHTTPBody(r, &c)

	if validateCarInventorySchema(c) {
		_, dbError := createCarInventory(c)
		if dbError != nil {
			responseDBError(w, r)
		} else {
			responseJSON(w, r, APIResponse{
				Success: true,
				Data:    c,
			})
		}
	} else {
		responseInvalidParameter(w, r)
	}
}

func handleDeleteInventoryCar(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	vin := ps.ByName("vin")

	_, dbError := deleteCarInventory(vin)
	if dbError != nil {
		responseDBError(w, r)
	} else {
		responseJSON(w, r, APIResponse{
			Success: true,
			Data:    vin,
		})
	}
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
