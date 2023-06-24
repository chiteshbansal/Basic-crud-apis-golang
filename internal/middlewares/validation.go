package middleware

import (
	model "first-api/internal/models"
	"first-api/pkg/validator"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ValidateCreateUser is a middleware function that validates user data from the request body.
func ValidateCreateUser(ctx *gin.Context) {
	// Bind JSON request body to struct
	var reqBody validator.CreateUser

	if err := ctx.BindJSON(&reqBody); err != nil {
		fmt.Println("Error: ", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Json binding Failed",
		})
		return
	}
	// Validate the request body
	if valError := reqBody.Validate(); valError != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": valError,
		})
		return
	}

	// Access the User and confirmPassword separately
	user := model.User{
		Name:     reqBody.Name,
		Email:    reqBody.Email,
		Password: reqBody.Password,
		Phone:    reqBody.Phone,
		Address:  reqBody.Address,
		Role:     reqBody.Role,
	}
	confirmPassword := reqBody.ConfirmPassword

	// Use the user and confirmPassword as needed in further middlewares

	ctx.Set("body", user)
	ctx.Set("confirmPassword", confirmPassword)
	ctx.Next()
}

// ValidateUpdateUser is a middleware function that validates user data from the request body.
func ValidateUpdateUser(ctx *gin.Context) {
	// Bind JSON request body to struct
	var reqBody validator.UpdateUser

	if err := ctx.BindJSON(&reqBody); err != nil {
		fmt.Println("Error: ", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Json binding Failed",
		})
		return
	}
	// Validate the request body
	if valError := reqBody.Validate(); valError != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": valError,
		})
		return
	}

	// Access the User and confirmPassword separately
	user := model.User{
		Name:    reqBody.Name,
		Email:   reqBody.Email,
		Phone:   reqBody.Phone,
		Address: reqBody.Address,
		Role:    reqBody.Role,
	}

	// Use the user and confirmPassword as needed in further middlewares

	ctx.Set("body", user)
	ctx.Next()
}

// ValidateLogin is a middleware function that validates user data from the request body.
func ValidateLogin(ctx *gin.Context) {
	// Bind JSON request body to struct
	var reqBody validator.Login

	if err := ctx.BindJSON(&reqBody); err != nil {
		fmt.Println("Error: ", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Json binding Failed",
		})
		return
	}
	// Validate the request body
	if valError := reqBody.Validate(); valError != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": valError,
		})
		return
	}

	// Access the User and confirmPassword separately
	user := model.User{
		Email:    reqBody.Email,
		Password: reqBody.Password,
	}

	// Use the user and confirmPassword as needed in further middlewares

	ctx.Set("body", user)
	ctx.Next()
}

// ValidateCreatePost is a middleware function that validates post data from the request body.
func ValidateCreatePost(ctx *gin.Context) {
	// Bind JSON request body to struct
	var reqBody validator.CreatePost

	if err := ctx.BindJSON(&reqBody); err != nil {
		fmt.Println("Error: ", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Json binding Failed",
		})
		return
	}
	// Validate the request body
	if valError := reqBody.Validate(); valError != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": valError,
		})
		return
	}

	// Access the Post separately
	post := model.Post{
		Title:       reqBody.Title,
		Description: reqBody.Description,
		Creator:     reqBody.Creator,
	}

	// Use the post as needed in further middlewares

	ctx.Set("body", post)
	ctx.Next()
}

// ValidateUpdatePost is a middleware function that validates post data from the request body.
func ValidateUpdatePost(ctx *gin.Context) {
	// Bind JSON request body to struct
	var reqBody validator.UpdatePost

	if err := ctx.BindJSON(&reqBody); err != nil {
		fmt.Println("Error: ", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Json binding Failed",
		})
		return
	}
	// Validate the request body
	if valError := reqBody.Validate(); valError != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": valError,
		})
		return
	}

	// Access the Post separately
	post := model.Post{
		Title:       reqBody.Title,
		Description: reqBody.Description,
		Creator:     reqBody.Creator,
		CreatorId:   reqBody.CreatorId,
	}

	// Use the post as needed in further middlewares

	ctx.Set("body", post)
	ctx.Next()
}

// AddComment is a middleware function that validates comment data from the request body.
func AddComment(ctx *gin.Context) {
	// Bind JSON request body to struct
	var reqBody struct {
		ID      string            `json:"id"`
		Comment validator.Comment `json:"comment"`
	}

	if err := ctx.BindJSON(&reqBody); err != nil {
		fmt.Println("Error:", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "JSON binding failed",
		})
		return
	}

	// Validate the request body
	if valErr := reqBody.Comment.Validate(); valErr != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": valErr.Error(),
		})
		return
	}

	ctx.Set("body", reqBody)
	ctx.Next()
}
