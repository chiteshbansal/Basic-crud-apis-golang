// This is the main package for the first-api application.
package main

import (
	"fmt"
	
	// Import custom packages
	model "first-api/api/Models"
	route "first-api/api/Routes"
	"first-api/api/repository"
	"first-api/api/service"
	"first-api/api/utils"
	db "first-api/database"

	// Import third party packages
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	_ "github.com/go-sql-driver/mysql" // This is a driver for MySQL to be used with the gorm package
)

var err error // Global error variable

func main() {
	// Set up environment variables
	viper.SetConfigFile("./.env") // specify location of the .env file
	viper.ReadInConfig() // read the .env file

	// Initialize user service with a repository
	userService := &service.UserService{
		Store: repository.UserStore{},
	}

	// Connect to the MySQL database using gorm
	db.DB, err = gorm.Open("mysql", db.DbURL(db.BuildConfig()))
	if err != nil {
		fmt.Println("Status:", err)
	}

	// Make sure to close the database connection when the main function exits
	defer db.DB.Close()

	// Automigrate user model, this will create the user table in the database
	db.DB.AutoMigrate(&model.User{})

	// Create a new gin engine
	r := gin.Default()
	
	// Register routes
	utils.RegisterRoutes(userService) // register user service related routes
	route.InitializeRoutes(r) // initialize other routes

	// Run the gin engine
	r.Run()
}
