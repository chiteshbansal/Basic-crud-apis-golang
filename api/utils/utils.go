package utils

import (
	route "first-api/api/Routes"
	"first-api/api/service"
)

func RegisterRoutes(userService *service.UserService) {

	route.RegisterRoutes(route.RouteDef{
		Path:    "/user",
		Version: "v1",
		Method:  "GET",
		Handler: userService.GetUsers,
	})
	route.RegisterRoutes(route.RouteDef{
		Path:    "/user",
		Version: "v1",
		Method:  "POST",
		Handler: userService.CreateUser,
	})

	route.RegisterRoutes(route.RouteDef{
		Path:    "/user/:id",
		Version: "v1",
		Method:  "PUT",
		Handler: userService.UpdateUser,
	})

	route.RegisterRoutes(route.RouteDef{
		Path:    "/user/:id",
		Version: "v1",
		Method:  "DELETE",
		Handler: userService.DeleteUser,
	})
	route.RegisterRoutes(route.RouteDef{
		Path:    "/user/filter",
		Version: "v1",
		Method:  "GET",
		Handler: userService.GetUser,
	})
}
