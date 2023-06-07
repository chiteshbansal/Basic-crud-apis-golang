package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

func GenerateJWT(email string) (string, error) {

	viper.SetConfigFile("../.env")
	viper.ReadInConfig()

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(2 * time.Hour).Unix()
	claims["authorized"] = true
	claims["user"] = email
	// fmt.Println([]byte(viper.GetString("SECRET_KEY")), ))
	tokenString, err := token.SignedString([]byte((viper.GetString("SECRET_KEY"))))
	if err != nil {
		return " ", err
	}

	return tokenString, nil
}
