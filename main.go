package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Employee struct {
	ID         int    `json:"employeeId"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Department string `json:"department"`
}

var employees = map[int]Employee{
	101: {
		ID:         101,
		Name:       "Siva",
		Email:      "siva@example.com",
		Department: "Platform Engineering",
	},
	102: {
		ID:         102,
		Name:       "Ravi",
		Email:      "ravi@example.com",
		Department: "Engineering",
	},
}

func setupRouter() *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	router.GET("/employees/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid employee ID",
			})
			return
		}

		employee, exists := employees[id]

		if !exists {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Employee not found",
			})
			return
		}

		c.JSON(http.StatusOK, employee)
	})

	return router
}

func main() {
	router := setupRouter()

	router.Run(":8080")
}
