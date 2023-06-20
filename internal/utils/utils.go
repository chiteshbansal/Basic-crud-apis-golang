package utils

import (
	"bytes"
	"encoding/gob"
	model "first-api/internal/models"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

func GenerateJWT(user *model.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(2 * time.Hour).Unix()
	claims["authorized"] = true
	claims["user"] = user.Email
	claims["role"] = user.Role
	claims["id"] = user.Id
	tokenString, err := token.SignedString([]byte((viper.GetString("SECRET_KEY"))))
	if err != nil {
		return " ", err
	}
	return tokenString, nil
}

func InterfaceToBytes(data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(data)
	if err != nil {
		return nil, fmt.Errorf("error encoding data: %v", err)
	}
	return buf.Bytes(), nil
}
