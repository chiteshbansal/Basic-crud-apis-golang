// Package service provides functionalities for user-related operations.
package service

import (
	"context"
	"encoding/json"
	"errors"
	model "first-api/internal/models"
	"first-api/internal/repository"
	route "first-api/internal/route"
	"first-api/internal/utils"
	"first-api/pkg/cache"
	"fmt"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

// User encapsulates use case logic for users.
type User struct {
	Store     repository.UserStorer
	UserCache cache.UserCache
}

// CreateUser creates a new user by hashing the password and storing the user in the database.
func (u *User) CreateUser(ctx context.Context, req *route.AppReq) route.AppResp {
	var user model.User
	jsonData, _ := json.Marshal(req.Body)
	json.Unmarshal(jsonData, &user)

	// Ensure the password and confirmation password match.
	if req.Body["password"].(string) != req.Body["confirmPassword"].(string) {
		return map[string]interface{}{
			"status": http.StatusBadRequest,
			"error":  "Password and confirm password do not match!",
		}
	}

	// Hash the password using bcrypt.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Body["password"].(string)), bcrypt.DefaultCost)
	if err != nil {
		return map[string]interface{}{
			"status": http.StatusInternalServerError,
			"error":  "Failed to hash password!",
		}
	}

	user.Password = string(hashedPassword)

	// Create the user in the database.
	err = u.Store.CreateUser(&user)
	if err != nil {
		return map[string]interface{}{
			"status": http.StatusNotFound,
			"error":  err.Error(),
		}
	}

	user.Password = "" // Clear the password for the response.
	return map[string]interface{}{
		"status": http.StatusOK,
		"user":   user,
	}
}

// GetUsers retrieves all users from the database.
func (u *User) GetUsers(ctx context.Context, req *route.AppReq) route.AppResp {
	var users []model.User
	err := u.Store.GetAllUsers(&users)
	if err != nil {
		return map[string]interface{}{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		}
	}
	return map[string]interface{}{
		"status": http.StatusOK,
		"users":  users,
	}
}

// UpdateUser updates a user in the database based on the given user ID.
func (u *User) UpdateUser(ctx context.Context, req *route.AppReq) route.AppResp {
	id := req.Params["id"]
	var user model.User
	err := u.Store.GetUser(&user, "id="+id)
	if err != nil {
		return map[string]interface{}{
			"status": http.StatusNotFound,
			"error":  err.Error(),
		}
	}

	jsonData, _ := json.Marshal(req.Body)
	json.Unmarshal(jsonData, &user)

	val, _ := strconv.ParseUint(id, 10, 64) // Convert string to uint.
	user.Id = uint(val)

	err = u.Store.UpdateUser(&user, id)
	if err != nil {
		return map[string]interface{}{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		}
	}
	return map[string]interface{}{
		"status":  http.StatusOK,
		"message": "User updated !!",
		"user":    user,
	}
}

// DeleteUser removes a user from the database based on the given user ID.
func (u *User) DeleteUser(ctx context.Context, req *route.AppReq) route.AppResp {
	var user model.User
	id := req.Params["id"]
	err := u.Store.GetUser(&user, "id="+id)
	if err != nil {
		return map[string]interface{}{
			"status": http.StatusNotFound,
			"error":  err.Error(),
		}
	}

	jsonData, _ := json.Marshal(req.Body)
	json.Unmarshal(jsonData, &user)

	err = u.Store.DeleteUser(&user, id)
	if err != nil {
		return map[string]interface{}{
			"status": http.StatusNotFound,
			"error":  err.Error(),
		}
	}
	return map[string]interface{}{
		"status":  http.StatusOK,
		"message": "User with " + id + " is Deleted!",
	}
}

// GetUser retrieves a user from the cache if present, else retrieves from the database.
func (u *User) GetUser(ctx context.Context, req *route.AppReq) route.AppResp {
	query := req.Query["filter"] + "=" + req.Query["value"]
	var user *model.User

	userInterface, _ := u.UserCache.Get(query)

	if userInterface == nil {
		fmt.Println("Not cached!!")
		user = &model.User{}
		err := u.Store.GetUser(user, query)
		if err != nil {
			return map[string]interface{}{
				"status": http.StatusInternalServerError,
				"error":  err.Error(),
			}
		}
		u.UserCache.Set(query, user, nil) // Set the user in the cache.
		return map[string]interface{}{
			"status": http.StatusOK,
			"user":   user,
		}
	} else {
		userInterface, _ := u.UserCache.Get(query)
		userMap, ok := userInterface.(map[string]interface{})
		if !ok {
			return map[string]interface{}{
				"status": http.StatusInternalServerError,
				"error":  errors.New("userInterface is not mapStringInterface"),
			}
		}

		// Now you can unmarshal userMap into your user struct.
		userBytes, err := json.Marshal(userMap)
		if err != nil {
			return map[string]interface{}{
				"status": http.StatusInternalServerError,
				"error":  errors.New("json Marshal failed!!"),
			}
		}

		var user *model.User
		err = json.Unmarshal(userBytes, &user)
		if err != nil {
			return map[string]interface{}{
				"status": http.StatusInternalServerError,
				"error":  errors.New("JSON unmarshal failed"),
			}
		}

		fmt.Println("using cached data")
		return map[string]interface{}{
			"status": http.StatusOK,
			"user":   user,
		}
	}
}

// Login verifies the user credentials, generates a JWT token and returns it.
func (u *User) Login(ctx context.Context, req *route.AppReq) route.AppResp {
	query := "email=\"" + req.Body["email"].(string) + "\""
	var user model.User

	err := u.Store.GetUser(&user, query)
	if err != nil {
		return map[string]interface{}{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Body["password"].(string)))
	if err != nil {
		return map[string]interface{}{
			"status": http.StatusBadRequest,
			"error":  "Email Id and password do not match",
		}
	}

	token, err := utils.GenerateJWT(user.Email)
	if err != nil {
		return map[string]interface{}{
			"status": http.StatusInternalServerError,
			"error":  "Something went wrong",
		}
	}
	return map[string]interface{}{
		"status": http.StatusOK,
		"token":  token,
		"email":  user.Email,
	}
}
