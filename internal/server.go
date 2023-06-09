package server

import (

	// Import custom packages

	middleware "first-api/internal/middlewares"
	"first-api/internal/repository"
	route "first-api/internal/route"
	"first-api/internal/service"
	cache "first-api/pkg/cache"

	// Import third party packages
	"github.com/gin-gonic/gin"
)

// RegisterRoutes function registers routes for the user service.
// RegisterRoutes function registers routes for the user service.
func RegisterRoutes() {
	// Register GET route to retrieve all users.
	userService := &service.User{
		Store:     &repository.UserStore{},
		UserCache: cache.NewRedisCache("localhost:6379", 0, 1000),
	}
	route.RegisterRoutes(route.RouteDef{
		Path:    "/user",
		Version: "v1",
		Method:  "GET",
		Handler: userService.GetUsers,
		Middlewares: []gin.HandlerFunc{
			func(ctx *gin.Context) {
				middleware.VerifyJWT(ctx, userService.UserCache)
			},
		},
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
		Path:    "/user/:id",
		Version: "v1",
		Method:  "PUT",
		Handler: userService.UpdateUser,
		Middlewares: []gin.HandlerFunc{func(ctx *gin.Context) {
			middleware.VerifyJWT(ctx, userService.UserCache)
		}, middleware.ValidateUserData},
	})

	// Register DELETE route to remove a user.
	route.RegisterRoutes(route.RouteDef{
		Path:    "/user/:id",
		Version: "v1",
		Method:  "DELETE",
		Handler: userService.DeleteUser,
		Middlewares: []gin.HandlerFunc{func(ctx *gin.Context) {
			middleware.VerifyJWT(ctx, userService.UserCache)
		}, middleware.ValidateUserData},
	})

	// Register GET route to retrieve a user based on filters.
	route.RegisterRoutes(route.RouteDef{
		Path:    "/user/filter",
		Version: "v1",
		Method:  "GET",
		Handler: userService.GetUser,
		Middlewares: []gin.HandlerFunc{func(ctx *gin.Context) {
			middleware.VerifyJWT(ctx, userService.UserCache)
		}},
	})
	route.RegisterRoutes(route.RouteDef{
		Path:    "/user/login",
		Version: "v1",
		Method:  "POST",
		Handler: userService.Login,
	})
}
