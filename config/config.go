package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"webserver.service/model"
)

type Config struct {
	Router *mux.Router
	DB     *sql.DB
}

func (config *Config) Initialize(DbUser string, DbPassword string, DbHost string, DbPort string, DbName string) error {
	connectionString := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v",
		DbUser, DbPassword, DbHost, DbPort, DbName)

	// Initialize the database connection here
	var err error
	config.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		return err
	}

	config.Router = mux.NewRouter().StrictSlash(true)
	config.handleRoutes()
	return nil
}

func (config *Config) Run() {
	log.Fatal(http.ListenAndServe(RouterAddr, config.Router))

}

func (config *Config) handleRoutes() {
	config.Router.HandleFunc("/products", config.getProducts).Methods("GET")
	config.Router.HandleFunc("/products/{id}", config.getProduct).Methods("GET")
	config.Router.HandleFunc("/product", config.createProduct).Methods("POST")
	config.Router.HandleFunc("/product/{id}", config.updateProduct).Methods("PUT")
	config.Router.HandleFunc("/product/{id}", config.deleteProduct).Methods("DELETE")
}

func SendResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

func sendError(w http.ResponseWriter, statusCode int, err string) {
	error_message := map[string]string{"error": err}
	SendResponse(w, statusCode, error_message)
}

func (config *Config) getProducts(w http.ResponseWriter, r *http.Request) {
	products, err := model.GetProducts(config.DB)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	SendResponse(w, http.StatusOK, products)
}

func (config *Config) getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}
	p := model.Product{ID: productID}
	err = p.GetProduct(config.DB)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			sendError(w, http.StatusNoContent, "Doesn't exist product ID")
		default:
			sendError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	SendResponse(w, http.StatusOK, p)
}

func (config *Config) createProduct(w http.ResponseWriter, r *http.Request) {
	var p model.Product
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := p.CreateProduct(config.DB); err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	SendResponse(w, http.StatusCreated, p)
}

func (config *Config) updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var p model.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	p.ID = productID

	err = p.UpdateProduct(config.DB)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	SendResponse(w, http.StatusOK, p)

}

func (config *Config) deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productId, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	p := model.Product{ID: productId}
	err = p.DeleteProduct(config.DB)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			sendError(w, http.StatusNotFound, "Doesn't exist product ID")
		default:
			sendError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	SendResponse(w, http.StatusNoContent, map[string]string{"message": "Successful deletion"})
}
