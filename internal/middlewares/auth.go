package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

func VerifyJWT(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	email := ctx.GetHeader("X-User-Email")
	viper.SetConfigFile("../.env")
	viper.ReadInConfig()

	if authHeader != "" {
		tokenString := strings.Split(authHeader, "Bearer ")[1]
		// Replace "YOUR_SECRET_KEY" with your actual secret key
		secretKey := []byte(viper.GetString("SECRET_KEY"))

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("Invalid signing method")
			}
			return secretKey, nil
		})

		if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		if token.Valid {
			// Verify the required claims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				ctx.AbortWithError(http.StatusUnauthorized, errors.New("Invalid token claims"))
				return
			}

			// Verify the user's identity or any other relevant information
			userEmail, ok := claims["user"].(string)
			if !ok || userEmail != email {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid User"})
				ctx.Abort()
				return
			}

			ctx.Next()
			return
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Token"})
			ctx.Abort()
			return
		}
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing Token"})
		ctx.Abort()
		return
	}
}
