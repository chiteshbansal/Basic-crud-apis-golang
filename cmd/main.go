// This is the main package for the first-api application.
package main

import (

	// Import custom packages
	server "first-api/internal"
	db "first-api/internal/database"
	route "first-api/internal/route"
	"log"

	// Import third party packages
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" // This is a driver for MySQL to be used with the gorm package
	"github.com/spf13/viper"
)

func init() {
	// Set up environment variables
	viper.SetConfigFile("../.env") // specify location of the .env file
	viper.ReadInConfig()           // read the .env file
}

func main() {
	// Initialize database
	err := db.NewDB()
	if err != nil {
		// Handle the error in an appropriate way, such as logging it or exiting the program
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Make sure to close the database connection when the main function exits
	defer db.DB.Close()

	// Create a new gin engine
	r := gin.Default()

	// Register routes
	server.RegisterRoutes(r) // register user service related routes

	route.InitializeRoutes(r) // initialize other routes

	// Run the gin engine
	r.Run()
}
