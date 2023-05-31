package service

import (
	"context"
	"encoding/json"
	model "first-api/api/Models"
	route "first-api/api/Routes"
	"first-api/api/repository"
	"fmt"
	"net/http"
	"strconv"
)

type UserService struct {
	Store repository.UserStore
}

func (u *UserService) CreateUser(ctx context.Context, req *route.AppReq) route.AppResp {
	var user model.User

	jsonData, err := json.Marshal(req.Body)
	json.Unmarshal(jsonData, &user)

	err = u.Store.CreateUser(ctx, &user)
	if err != nil {
		return map[string]interface{}{
			"status": http.StatusNotFound,
			"error":  err.Error(),
		}
	} else {
		return map[string]interface{}{
			"status": http.StatusOK,
			"user":   user,
		}
	}
}

func (u *UserService) GetUsers(ctx context.Context, req *route.AppReq) route.AppResp {
	var users []model.User
	err := u.Store.GetAllUsers(&users)
	if err != nil {
		fmt.Println(err.Error())
		return map[string]interface{}{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		}
	} else {
		return map[string]interface{}{
			"status": http.StatusOK,
			"users":  users,
		}
	}

}

// // update user data
func (u *UserService) UpdateUser(ctx context.Context, req *route.AppReq) route.AppResp {
	id := req.Params["id"]
	var user model.User
	query := "id=" + id

	err := u.Store.GetUser(&user, query)
	if err != nil {
		fmt.Println(err.Error())
		return map[string]interface{}{
			"status": http.StatusNotFound,
			"error":  err.Error(),
		}
	}

	jsonData, err := json.Marshal(req.Body)
	json.Unmarshal(jsonData, &user)

	// parse string to uint
	val, _ := (strconv.ParseUint(id, 10, 64))
	user.Id = (uint)(val)

	u.Store.UpdateUser(&user, id)

	if err != nil {
		return map[string]interface{}{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		}
	} else {
		return map[string]interface{}{
			"status":  http.StatusOK,
			"message": "User updated !!",
			"user":    user,
		}
	}

}

// // delete user

func (u *UserService) DeleteUser(ctx context.Context, req *route.AppReq) route.AppResp {
	var user model.User
	id := req.Params["id"]
	query := "id=" + id
	err := u.Store.GetUser(&user, query)
	if err != nil {
		fmt.Println(err.Error())
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
	} else {
		return map[string]interface{}{
			"status":  http.StatusOK,
			"message": "User with " + id + " is Deleted!",
		}
	}
}

func (u *UserService) GetUser(ctx context.Context, req *route.AppReq) route.AppResp {
	query := req.Query["filter"] + "=" + req.Query["value"]
	var user model.User

	err := u.Store.GetUser(&user, query)
	if err != nil {
		return map[string]interface{}{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		}
	} else {
		return map[string]interface{}{
			"status": http.StatusOK,
			"user":   user,
		}
	}
}
