package database

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var err error

type Handler struct {
	DB *sql.DB
	// Err error
}

type Customer struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

func (h *Handler) CreateTableHandler() {
	stmt, err := h.DB.Prepare(`CREATE TABLE IF NOT EXISTS customers(
							id SERIAL PRIMARY KEY,
							name TEXT,
							email TEXT,
							status TEXT
						);`)
	if _, err = stmt.Exec(); err != nil {
		log.Fatal(err)
		return
	}
}

func (h *Handler) CreateCustomerHandler(c *gin.Context) {
	cus := Customer{}
	if err := c.ShouldBindJSON(&cus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	stmt, err := h.DB.Prepare(`INSERT INTO customers(name, email, status)
							values ($1, $2, $3)
							RETURNING id`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	row := stmt.QueryRow(cus.Name, cus.Email, cus.Status)
	err = row.Scan(&cus.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cus)
}

func (h *Handler) GetCustomerHandler(c *gin.Context) {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	stmt, err := h.DB.Prepare(`SELECT id, name, email, status 
							FROM customers
							WHERE id=$1`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	cus := Customer{}
	row := stmt.QueryRow(ID)
	err = row.Scan(&cus.ID, &cus.Name, &cus.Email, &cus.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cus)
}

func (h *Handler) GetAllCustomerHandler(c *gin.Context) {
	stmt, err := h.DB.Prepare(`SELECT id, name, email, status 
							FROM customers`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	cus := []Customer{}
	row, err := stmt.Query()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	defer row.Close()

	var cusID int
	var name, email, status string
	for row.Next() {
		err = row.Scan(&cusID, &name, &email, &status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
		cus = append(cus, Customer{
			ID:     cusID,
			Name:   name,
			Email:  email,
			Status: status,
		})
	}
	c.JSON(http.StatusOK, cus)
}

func (h *Handler) UpdateCustomerHandler(c *gin.Context) {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	cus := Customer{ID: ID}
	if err := c.ShouldBindJSON(&cus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	stmt, err := h.DB.Prepare(`UPDATE customers
							SET name=$2, email=$3, status=$4
							WHERE id=$1`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	if _, err := stmt.Exec(cus.ID, cus.Name, cus.Email, cus.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cus)
}

func (h *Handler) DeleteCustomerHandler(c *gin.Context) {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	stmt, err := h.DB.Prepare(`DELETE FROM customers
							WHERE id=$1`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	if _, err := stmt.Exec(ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "customer deleted"})

}
