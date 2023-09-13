package main

import (
	"os"
	"server/Connections"
	"github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "5000"
	}

	router := gin.New() 
	router.Use(gin.Logger()) // this is a middleware that logs the requests

	router.Use(cors.Default()) // this is a middleware that allows cross-origin requests
							   // Allow all origin, allow all CRUD operations, allow all headers
	// these are the endpoints
	//C
	router.POST("/order/create", routes.CreateOrder)
	//R
	router.GET("/waiter/:server", routes.GetOrdersByServer)
	router.GET("/orders", routes.GetOrders)
	router.GET("/order/:id/", routes.GetOrderById)
	//U
	router.PUT("/waiter/update/:id", routes.UpdateServer)
	router.PUT("/order/update/:id", routes.UpdateOrder)
	//D
	router.DELETE("/order/delete/:id", routes.DeleteOrder)

	//this runs the server and allows it to listen to requests.
	router.Run(":" + port) // listen on localhost and serve on port 5000
}