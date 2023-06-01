package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

func verifyJWT(ctx *gin.Context) {

	viper.SetConfigFile("../.env")
	viper.ReadInConfig()


	authHeader := ctx.GetHeader("Authorization")
	if authHeader != "" {
		tokenString := strings.Split(authHeader, "Bearer ")[1]

		// Replace "YOUR_SECRET_KEY" with your actual secret key
		secretKey := []byte(viper.GetString("SECRET_KEY"))

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return secretKey, nil
		})

		if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		if token.Valid {
			ctx.Next()
			return
		} else {
			ctx.AbortWithError(http.StatusUnauthorized, errors.New("Invalid token"))
			return
		}
	} else {
		ctx.AbortWithError(http.StatusUnauthorized, errors.New("Missing token"))
		return
	}
}
