package middleware

import (
	model "first-api/api/Models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

func ValidateUserData(ctx *gin.Context) {
	body := model.User{}

	if err := ctx.BindJSON((&body)); err != nil {
		fmt.Println("ERror ", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Something went wrong !!",
		})
		return
	}
	valError := validation.ValidateStruct(&body, validation.Field(&body.Name, validation.Required, validation.Length(5, 20)),
		validation.Field(&body.Email, validation.Required, is.Email),
		validation.Field(&body.Phone, validation.Required, validation.Length(10, 10)),
		validation.Field(&body.Address, validation.Required, validation.Length(10, 50)),
	)

	if valError != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": valError,
		})
		return
	}

	ctx.Set("body", body)
	ctx.Next()
}
