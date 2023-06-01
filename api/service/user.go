package service

import (
	"context"
	"encoding/json"
	"first-api/api/Models"
	"first-api/api/Routes"
	"first-api/api/repository"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

// UserService encapsulates use case logic for users.
type UserService struct {
	Store repository.UserStore
}

// CreateUser creates a new user.
func (u *UserService) CreateUser(ctx context.Context, req *route.AppReq) route.AppResp {
	var user model.User

	jsonData, err := json.Marshal(req.Body)
	json.Unmarshal(jsonData, &user)

	password := req.Body["password"].(string)
	confirmPassword := req.Body["confirmPassword"].(string)

	if password != confirmPassword {
		return map[string]interface{}{
			"status": http.StatusBadRequest,
			"error":  "Password and confirm password do not match!",
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return map[string]interface{}{
			"status": http.StatusInternalServerError,
			"error":  "Failed to hash password!",
		}
	}

	user.Password = string(hashedPassword)

	err = u.Store.CreateUser(&user)
	if err != nil {
		return map[string]interface{}{
			"status": http.StatusNotFound,
			"error":  err.Error(),
		}
	}

	// Clear the password field in the response for security purposes
	user.Password = ""

	return map[string]interface{}{
		"status": http.StatusOK,
		"user":   user,
	}
}


// GetUsers retrieves all users.
func (u *UserService) GetUsers(ctx context.Context, req *route.AppReq) route.AppResp {
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

// UpdateUser updates user data based on ID.
func (u *UserService) UpdateUser(ctx context.Context, req *route.AppReq) route.AppResp {
	id := req.Params["id"]
	var user model.User
	query := "id=" + id

	err := u.Store.GetUser(&user, query)
	if err != nil {
		return map[string]interface{}{
			"status": http.StatusNotFound,
			"error":  err.Error(),
		}
	}

	jsonData, err := json.Marshal(req.Body)
	json.Unmarshal(jsonData, &user)

	// parse string to uint
	val, _ := strconv.ParseUint(id, 10, 64)
	user.Id = uint(val)

	u.Store.UpdateUser(&user, id)

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

// DeleteUser removes a user based on ID.
func (u *UserService) DeleteUser(ctx context.Context, req *route.AppReq) route.AppResp {
	var user model.User
	id := req.Params["id"]
	query := "id=" + id
	err := u.Store.GetUser(&user, query)
	if err != nil {
		return map[string]interface{}{
			"status": http.StatusNotFound,
			"error":  err.Error(),
		}
	}

	jsonData, err := json.Marshal(req.Body)
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

// GetUser retrieves a user based on filter query.
func (u *UserService) GetUser(ctx context.Context, req *route.AppReq) route.AppResp {
	query := req.Query["filter"] + "=" + req.Query["value"]
	var user model.User

	err := u.Store.GetUser(&user, query)
	if err != nil {
		return map[string]interface{}{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		}
	}

	return map[string]interface{}{
		"status": http.StatusOK,
		"user":   user,
	}
}
