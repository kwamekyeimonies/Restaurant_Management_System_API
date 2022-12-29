package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kwamekyeimonies/restaurant_management_system_backend/database"

	// "github.com/kwamekyeimonies/restaurant_management_system_backend/middlewares"
	"github.com/kwamekyeimonies/restaurant_management_system_backend/routes"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.New()
	router.Use(gin.Logger())
	// routes.Use(middlewares.Authentication())

	routes.FoodRoutes(router)
	routes.InvoiceRoutes(router)
	routes.MenuRoutes(router)
	routes.OrderItemRoutes(router)
	routes.OrderRoutes(router)
	routes.TableRoutes(router)
	routes.UserRoutes(router)

	router.Run(":" + port)

}
