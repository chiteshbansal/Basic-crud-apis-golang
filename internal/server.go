package server

import (

	// Import custom packages

	middleware "first-api/internal/middlewares"
	"first-api/internal/ratelimiter"
	"first-api/internal/repository"
	route "first-api/internal/route"
	"first-api/internal/service"
	cache "first-api/pkg/cache"
	"time"

	// Import third party packages
	"github.com/gin-gonic/gin"
)

// RegisterRoutes function registers routes for the user service.
// RegisterRoutes function registers routes for the user service.
func RegisterRoutes(server *gin.Engine) {

	userCache := cache.GetRateLimiterCache()

	rateLimiterCache := cache.GetRateLimiterCache()

	tb := ratelimiter.NewTokenBucket(100, 1, time.Now())
	server.Use(func(ctx *gin.Context) {
		middleware.RateLimitMiddleware(ctx, rateLimiterCache, tb)
	})

	userService := &service.User{
		Store:     &repository.UserStore{},
		UserCache: userCache,
	}

	postService := &service.Post{
		Store:     &repository.PostStore{},
		UserCache: userCache,
	}

	// Register GET route to retrieve all users.
	route.RegisterRoutes(route.RouteDef{
		Path:    "/user",
		Version: "v1",
		Method:  "GET",
		Handler: userService.GetUsers,
		Middlewares: []gin.HandlerFunc{
			func(ctx *gin.Context) {
				middleware.VerifyJWT(ctx, userService.UserCache, "user")
			},
		},
	})

	// Register POST route to create a new user. Includes a middleware to validate user data.
	route.RegisterRoutes(route.RouteDef{
		Path:        "/user",
		Version:     "v1",
		Method:      "POST",
		Handler:     userService.CreateUser,
		Middlewares: []gin.HandlerFunc{middleware.ValidateCreateUser},
	})

	// Register PUT route to update a user. Includes a middleware to validate user data.
	route.RegisterRoutes(route.RouteDef{
		Path:    "/user/:id",
		Version: "v1",
		Method:  "PUT",
		Handler: userService.UpdateUser,
		Middlewares: []gin.HandlerFunc{func(ctx *gin.Context) {
			middleware.VerifyJWT(ctx, userService.UserCache, "admin")
		}, middleware.ValidateUpdateUser},
	})

	// Register DELETE route to remove a user.
	route.RegisterRoutes(route.RouteDef{
		Path:    "/user/:id",
		Version: "v1",
		Method:  "DELETE",
		Handler: userService.DeleteUser,
		Middlewares: []gin.HandlerFunc{func(ctx *gin.Context) {
			middleware.VerifyJWT(ctx, userService.UserCache, "admin")
		}},
	})

	// Register GET route to retrieve a user based on filters.
	route.RegisterRoutes(route.RouteDef{
		Path:    "/user/filter",
		Version: "v1",
		Method:  "GET",
		Handler: userService.GetUser,
		Middlewares: []gin.HandlerFunc{func(ctx *gin.Context) {
			middleware.VerifyJWT(ctx, userService.UserCache, "user")
		}},
	})
	route.RegisterRoutes(route.RouteDef{
		Path:    "/user/login",
		Version: "v1",
		Method:  "POST",
		Handler: userService.Login,
		Middlewares: []gin.HandlerFunc{middleware.ValidateLogin},
	})

	// Post routes

	// create post route
	route.RegisterRoutes(route.RouteDef{
		Path:    "/post",
		Version: "v1",
		Method:  "POST",
		Handler: postService.CreatePost,
		Middlewares: []gin.HandlerFunc{
			func(ctx *gin.Context) {
				middleware.VerifyJWT(ctx, userService.UserCache, "user")
			},
		},
	})

	// Get all post route
	route.RegisterRoutes(route.RouteDef{
		Path:    "/post",
		Version: "v1",
		Method:  "GET",
		Handler: postService.GetPosts,
		Middlewares: []gin.HandlerFunc{
			func(ctx *gin.Context) {
				middleware.VerifyJWT(ctx, userService.UserCache, "user")
			},
		},
	})
	// get post by id route
	route.RegisterRoutes(route.RouteDef{
		Path:    "/post/:id",
		Version: "v1",
		Method:  "GET",
		Handler: postService.GetPost,
		Middlewares: []gin.HandlerFunc{
			func(ctx *gin.Context) {
				middleware.VerifyJWT(ctx, userService.UserCache, "user")
			},
		},
	})

	// delete post route
	route.RegisterRoutes(route.RouteDef{
		Path:    "/post/:id",
		Version: "v1",
		Method:  "DELETE",
		Handler: postService.DeletePost,
		Middlewares: []gin.HandlerFunc{func(ctx *gin.Context) {
			middleware.VerifyJWT(ctx, postService.UserCache, "user")
		}},
	})

	// update post route
	// Register PUT route to update a user. Includes a middleware to validate user data.
	route.RegisterRoutes(route.RouteDef{
		Path:    "/post/:id",
		Version: "v1",
		Method:  "PUT",
		Handler: postService.UpdatePost,
		Middlewares: []gin.HandlerFunc{func(ctx *gin.Context) {
			middleware.VerifyJWT(ctx, postService.UserCache, "user")
		}},
	})

	// comment on post
	route.RegisterRoutes(route.RouteDef{
		Path:    "/post/comment/",
		Version: "v1",
		Method:  "POST",
		Handler: postService.AddComment,
		Middlewares: []gin.HandlerFunc{func(ctx *gin.Context) {
			middleware.VerifyJWT(ctx, postService.UserCache, "admin")
		}},
	})
}
