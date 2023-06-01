package utils

import (
	"time"

	// "github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)
func GenerateJWT() (string , error){

	viper.SetConfigFile("../.env")
	viper.ReadInConfig()

	token := jwt.New(jwt.SigningMethodEdDSA)

	claims :=token.Claims.(jwt.MapClaims)
	claims["exp"]  = time.Now().Add(2*time.Hour)
	claims["authorized"] = true
	claims["user"] = "username"

	tokenString,err :=token.SignedString([]byte(viper.GetString("SECRET_KEY")))
	if err!=nil{
		return " " ,err
	}

	// ctx.Writer.Header().Set("Authorization", "Bearer "+tokenString)
	
	return tokenString,nil
}