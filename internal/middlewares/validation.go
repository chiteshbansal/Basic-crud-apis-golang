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
	var reqBody validator.CreateUser

	if err := ctx.BindJSON(&reqBody); err != nil {
		fmt.Println("Error: ", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Json binding Failed",
		})
		return
	}
	if valError := reqBody.Validate(); valError != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": valError,
		})
		return
	}
	fmt.Println(reqBody)
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

// ValidateCreateUser is a middleware function that validates user data from the request body.
func ValidateUpdateUser(ctx *gin.Context) {
	var reqBody validator.UpdateUser

	if err := ctx.BindJSON(&reqBody); err != nil {
		fmt.Println("Error: ", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Json binding Failed",
		})
		return
	}
	if valError := reqBody.Validate(); valError != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": valError,
		})
		return
	}
	fmt.Println(reqBody)
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

// ValidateCreateUser is a middleware function that validates user data from the request body.
func ValidateLogin(ctx *gin.Context) {
	var reqBody validator.Login

	if err := ctx.BindJSON(&reqBody); err != nil {
		fmt.Println("Error: ", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Json binding Failed",
		})
		return
	}
	if valError := reqBody.Validate(); valError != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": valError,
		})
		return
	}
	fmt.Println(reqBody)
	// Access the User and confirmPassword separately
	user := model.User{
		Email:    reqBody.Email,
		Password: reqBody.Password,
	}
	// Use the user and confirmPassword as needed in further middlewares

	ctx.Set("body", user)
	ctx.Next()
}
