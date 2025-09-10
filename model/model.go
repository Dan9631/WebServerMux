package model

import (
	"database/sql"
	"fmt"
)

type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func GetProducts(db *sql.DB) ([]Product, error) {
	query := "SELECT id, name, quantity, price FROM products"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Quantity, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (p *Product) GetProduct(db *sql.DB) error {
	query := "SELECT id, name, quantity, price FROM products WHERE id = ?"
	row := db.QueryRow(query, p.ID)

	if err := row.Scan(&p.ID, &p.Name, &p.Quantity, &p.Price); err != nil {
		return err
	}
	return nil
}

func (p *Product) CreateProduct(db *sql.DB) error {
	query := fmt.Sprintf("INSERT INTO products(name,quantity,price) values('%v',%v,%v)", p.Name, p.Quantity, p.Price)
	result, err := db.Exec(query)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	p.ID = int(id)
	return nil
}

func (p *Product) UpdateProduct(db *sql.DB) error {
	query := fmt.Sprintf("UPDATE products SET name='%v', quantity=%v, price=%v WHERE id=%v", p.Name, p.Quantity, p.Price, p.ID)
	result, err := db.Exec(query)
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return err
}

func (p *Product) DeleteProduct(db *sql.DB) error {
	query := "DELETE FROM products WHERE id = ?"
	result, err := db.Exec(query, p.ID)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
