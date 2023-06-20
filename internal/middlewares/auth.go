package middleware

import (
	"errors"
	"first-api/pkg/cache"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

func VerifyJWT(ctx *gin.Context, cache cache.UserCache, role string) {
	authHeader := ctx.GetHeader("Authorization")
	email := ctx.GetHeader("X-User-Email")

	viper.SetConfigFile("../.env")
	viper.ReadInConfig()

	if authHeader != "" {
		tokenString := strings.Split(authHeader, "Bearer ")[1]
		// Replace "YOUR_SECRET_KEY" with your actual secret key
		secretKey := []byte(viper.GetString("SECRET_KEY"))

		// Check if the token is present in the cache.
		if user, _ := cache.Get(tokenString); user != nil {
			// If the token is found in the cache, pass the request to the next middleware function.
			ctx.Next()
			return
		}

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

			userRole, ok := claims["role"].(string)
			if !ok || role == "admin" && userRole != "admin" {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "UnAuthorized Feature access !!"})
				ctx.Abort()
				return
			}

			// If the token is valid, store it in the cache.
			// NOTE: You need to store some data instead of `nil` in the cache. This could be the user's ID or some other data.
			exp := claims["exp"].(float64)
			expTime := time.Unix(int64(exp), 0)

			// // Calculate the remaining time until the token expires
			remainingTime := time.Until(expTime)

			cache.Set(tokenString, "true", &remainingTime)

			userId, _ := claims["id"]
			ctx.Set("userId", userId)
			ctx.Set("role", userRole)

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
