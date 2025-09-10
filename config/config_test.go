package config

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupTestConfig() *Config {
	cfg := &Config{}
	err := cfg.Initialize("root", "rootpass", "localhost", "3306", "test") // Asume que la base de datos de prueba está configurada
	if err != nil {
		panic(err)
	}
	return cfg
}

func createTable(cfg *Config) {
	createTableQuery := `
	   create table IF NOT EXISTS  
	    products(id int NOT NULL AUTO_INCREMENT,
    	name varchar(255) NOT NULL,
    	quantity int,
    	price float(10,7),
    	PRIMARY KEY(id)
    );`
	_, err := cfg.DB.Exec(createTableQuery)
	if err != nil {
		panic(err)
	}
}

func cleanTable(cfg *Config) {
	cleanTableQuery := `DELETE FROM products;`
	_, err := cfg.DB.Exec(cleanTableQuery)
	if err != nil {
		panic(err)
	}
	alterIndexQuery := `ALTER TABLE products AUTO_INCREMENT = 1;`
	_, err = cfg.DB.Exec(alterIndexQuery)
	if err != nil {
		panic(err)
	}
}

func insertData(cfg *Config) {
	insertDataQuery := `
		INSERT INTO products (name, quantity, price) VALUES
		('Producto 1', 10, 99.99),
		('Producto 2', 20, 199.99),
		('Producto 3', 30, 299.99);`
	_, err := cfg.DB.Exec(insertDataQuery)
	if err != nil {
		panic(err)
	}
}

func TestGetProducts_OK(t *testing.T) {
	cfg := setupTestConfig()
	createTable(cfg)
	insertData(cfg)

	req := httptest.NewRequest("GET", "/products", nil)
	w := httptest.NewRecorder()
	cfg.Router.ServeHTTP(w, req)
	cleanTable(cfg)

	if w.Code != http.StatusOK {
		t.Errorf("Esperado status 200, obtenido %d", w.Code)
	}

	var products []map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &products); err != nil {
		t.Errorf("Error al parsear respuesta: %v", err)
	}
}

func TestGetProduct_OK(t *testing.T) {
	cfg := setupTestConfig()
	createTable(cfg)
	insertData(cfg)

	// Asume que el producto con ID 1 existe
	req := httptest.NewRequest("GET", "/products/1", nil)
	w := httptest.NewRecorder()
	cfg.Router.ServeHTTP(w, req)
	cleanTable(cfg)

	if w.Code != http.StatusOK {
		t.Errorf("Esperado status 200, obtenido %d", w.Code)
	}

	var product map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &product); err != nil {
		t.Errorf("Error al parsear respuesta: %v", err)
	}
}

func TestGetProduct_NotFound(t *testing.T) {
	cfg := setupTestConfig()
	createTable(cfg)
	insertData(cfg)
	// Asume que el producto con ID 999 no existe
	req := httptest.NewRequest("GET", "/products/999", nil)
	w := httptest.NewRecorder()
	cfg.Router.ServeHTTP(w, req)
	cleanTable(cfg)

	if w.Code != http.StatusNoContent {
		t.Errorf("Esperado status 204, obtenido %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Errorf("Error al parsear respuesta: %v", err)
	}
	if resp["error"] != "Doesn't exist product ID" {
		t.Errorf("Mensaje de error inesperado: %v", resp["error"])
	}
}

func TestGetProduct_BadRequest(t *testing.T) {
	cfg := setupTestConfig()
	req := httptest.NewRequest("GET", "/products/abc", nil)
	w := httptest.NewRecorder()
	cfg.Router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Esperado status 400, obtenido %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Errorf("Error al parsear respuesta: %v", err)
	}
	if resp["error"] != "Invalid product ID" {
		t.Errorf("Mensaje de error inesperado: %v", resp["error"])
	}
}

func TestCreateProduct_OK(t *testing.T) {
	cfg := setupTestConfig()
	body := `{"name":"Nuevo","quantity":10,"price":99.99}`
	req := httptest.NewRequest("POST", "/product", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	cfg.Router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Esperado status 201, obtenido %d", w.Code)
	}
	var product map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &product); err != nil {
		t.Errorf("Error al parsear respuesta: %v", err)
	}
}

func TestUpdateProduct_OK(t *testing.T) {
	cfg := setupTestConfig()
	createTable(cfg)
	insertData(cfg)

	body := `{"name":"Actualizado","quantity":20,"price":199.99}`
	req := httptest.NewRequest("PUT", "/product/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	cfg.Router.ServeHTTP(w, req)
	cleanTable(cfg)
	if w.Code != http.StatusOK {
		t.Errorf("Esperado status 200, obtenido %d, cuerpo: %s", w.Code, w.Body.String())
	}
	var product map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &product); err != nil {
		t.Errorf("Error al parsear respuesta: %v", err)
	}
}

func TestDeleteProduct_OK(t *testing.T) {
	cfg := setupTestConfig()
	createTable(cfg)
	insertData(cfg)

	req := httptest.NewRequest("DELETE", "/product/1", nil)
	w := httptest.NewRecorder()
	cfg.Router.ServeHTTP(w, req)
	cleanTable(cfg)

	if w.Code != http.StatusNoContent {
		t.Errorf("Esperado status 204, obtenido %d", w.Code)
	}
}
