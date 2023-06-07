// This is the main package for the first-api application.
package main

import (
	"fmt"

	// Import custom packages
	db "first-api/internal/database"
	middleware "first-api/internal/middlewares"
	model "first-api/internal/models"
	"first-api/internal/repository"
	route "first-api/internal/routes"
	"first-api/internal/service"
	"first-api/pkg/cache"

	// Import third party packages
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" // This is a driver for MySQL to be used with the gorm package
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

var err error // Global error variable

func main() {
	// Set up environment variables
	viper.SetConfigFile("./.env") // specify location of the .env file
	viper.ReadInConfig()          // read the .env file

	// Initialize user service with a repository
	userService := service.NewUser(
		&repository.UserStore{},
		cache.NewRedisCache("localhost:6379", 0, 10),
	)

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
	Registerroutes(userService) // register user service related routes
	route.Initializeroutes(r)   // initialize other routes

	// Run the gin engine
	r.Run()
}

// Registerroutes function registers routes for the user service.
func Registerroutes(userService *service.User) {
	// Register GET route to retrieve all users.
	route.Registerroutes(route.RouteDef{
		Path:        "/user",
		Version:     "v1",
		Method:      "GET",
		Handler:     userService.GetUsers,
		Middlewares: []gin.HandlerFunc{middleware.VerifyJWT},
	})

	// Register POST route to create a new user. Includes a middleware to validate user data.
	route.Registerroutes(route.RouteDef{
		Path:        "/user",
		Version:     "v1",
		Method:      "POST",
		Handler:     userService.CreateUser,
		Middlewares: []gin.HandlerFunc{middleware.ValidateUserData},
	})

	// Register PUT route to update a user. Includes a middleware to validate user data.
	route.Registerroutes(route.RouteDef{
		Path:    "/user/:id",
		Version: "v1",
		Method:  "PUT",
		Handler: userService.UpdateUser,
		Middlewares: []gin.HandlerFunc{middleware.
			VerifyJWT, middleware.ValidateUserData},
	})

	// Register DELETE route to remove a user.
	route.Registerroutes(route.RouteDef{
		Path:    "/user/:id",
		Version: "v1",
		Method:  "DELETE",
		Handler: userService.DeleteUser,
		Middlewares: []gin.HandlerFunc{middleware.
			VerifyJWT, middleware.ValidateUserData},
	})

	// Register GET route to retrieve a user based on filters.
	route.Registerroutes(route.RouteDef{
		Path:    "/user/filter",
		Version: "v1",
		Method:  "GET",
		Handler: userService.GetUser,
		Middlewares: []gin.HandlerFunc{middleware.
			VerifyJWT},
	})
	route.Registerroutes(route.RouteDef{
		Path:    "/user/login",
		Version: "v1",
		Method:  "POST",
		Handler: userService.Login,
	})
}