package utils

import (
	route "first-api/api/Routes"
	middleware "first-api/api/middlewares"
	"first-api/api/service"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(userService *service.UserService) {

	route.RegisterRoutes(route.RouteDef{
		Path:    "/user",
		Version: "v1",
		Method:  "GET",
		Handler: userService.GetUsers,
		// Middlewares: []gin.HandlerFunc{middleware.ValidateUserData},
	})
	route.RegisterRoutes(route.RouteDef{
		Path:        "/user",
		Version:     "v1",
		Method:      "POST",
		Handler:     userService.CreateUser,
		Middlewares: []gin.HandlerFunc{middleware.ValidateUserData},
	})

	route.RegisterRoutes(route.RouteDef{
		Path:        "/user/:id",
		Version:     "v1",
		Method:      "PUT",
		Handler:     userService.UpdateUser,
		Middlewares: []gin.HandlerFunc{middleware.ValidateUserData},
	})

	route.RegisterRoutes(route.RouteDef{
		Path:    "/user/:id",
		Version: "v1",
		Method:  "DELETE",
		Handler: userService.DeleteUser,
		// Middlewares: []gin.HandlerFunc{middleware.ValidateUserData},
	})
	route.RegisterRoutes(route.RouteDef{
		Path:    "/user/filter",
		Version: "v1",
		Method:  "GET",
		Handler: userService.GetUser,
		// Middlewares: []gin.HandlerFunc{middleware.ValidateUserData},
	})
}
