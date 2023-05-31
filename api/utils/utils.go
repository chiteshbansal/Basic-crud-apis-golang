// The utils package provides utility functions for the first-api application.
package utils

import (
	// Import custom packages
	route "first-api/api/Routes"
	middleware "first-api/api/middlewares"
	"first-api/api/service"

	// Import third party package
	"github.com/gin-gonic/gin"
)

// RegisterRoutes function registers routes for the user service.
func RegisterRoutes(userService *service.UserService) {
	// Register GET route to retrieve all users.
	route.RegisterRoutes(route.RouteDef{
		Path:    "/user",
		Version: "v1",
		Method:  "GET",
		Handler: userService.GetUsers,
		// Middlewares: []gin.HandlerFunc{middleware.ValidateUserData},
	})

	// Register POST route to create a new user. Includes a middleware to validate user data.
	route.RegisterRoutes(route.RouteDef{
		Path:        "/user",
		Version:     "v1",
		Method:      "POST",
		Handler:     userService.CreateUser,
		Middlewares: []gin.HandlerFunc{middleware.ValidateUserData},
	})

	// Register PUT route to update a user. Includes a middleware to validate user data.
	route.RegisterRoutes(route.RouteDef{
		Path:        "/user/:id",
		Version:     "v1",
		Method:      "PUT",
		Handler:     userService.UpdateUser,
		Middlewares: []gin.HandlerFunc{middleware.ValidateUserData},
	})

	// Register DELETE route to remove a user.
	route.RegisterRoutes(route.RouteDef{
		Path:    "/user/:id",
		Version: "v1",
		Method:  "DELETE",
		Handler: userService.DeleteUser,
		// Middlewares: []gin.HandlerFunc{middleware.ValidateUserData},
	})

	// Register GET route to retrieve a user based on filters.
	route.RegisterRoutes(route.RouteDef{
		Path:    "/user/filter",
		Version: "v1",
		Method:  "GET",
		Handler: userService.GetUser,
		// Middlewares: []gin.HandlerFunc{middleware.ValidateUserData},
	})
}
