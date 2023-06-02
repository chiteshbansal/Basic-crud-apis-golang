package middleware

import (
	model "first-api/api/Models"
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// CreateUserRequest is a struct that includes both User and confirmPassword fields.
type CreateUserRequest struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Address         string `json:"address"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

// ValidateUserData is a middleware function that validates user data from the request body.
func ValidateUserData(ctx *gin.Context) {
	reqBody := CreateUserRequest{}

	if err := ctx.BindJSON(&reqBody); err != nil {
		fmt.Println("Error: ", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Something went wrong!",
		})
		return
	}
	valError := validation.ValidateStruct(&reqBody,
		validation.Field(&reqBody.Name, validation.Required, validation.Length(5, 20)),
		validation.Field(&reqBody.Email, validation.Required, is.Email),
		validation.Field(&reqBody.Phone, validation.Required, validation.Length(10, 10)),
		validation.Field(&reqBody.Address, validation.Required, validation.Length(10, 50)),
	)

	if valError != nil {
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
	}
	confirmPassword := reqBody.ConfirmPassword

	// Use the user and confirmPassword as needed in further middlewares

	ctx.Set("body", user)
	ctx.Set("confirmPassword", confirmPassword)
	ctx.Next()
}
