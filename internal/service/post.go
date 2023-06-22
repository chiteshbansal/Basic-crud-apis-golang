// Package service provides functionalities for user-related operations.
package service

import (
	"context"
	"encoding/json"
	"errors"
	model "first-api/internal/models"
	"first-api/internal/repository"
	route "first-api/internal/route"
	"first-api/pkg/cache"
	"net/http"
	// "strconv"
)

// User encapsulates use case logic for users.
type Post struct {
	Store     repository.PostStorer
	UserCache cache.UserCache
}

// CreateUser creates a new user by hashing the password and storing the user in the database.
func (u *Post) CreatePost(ctx context.Context, req *route.AppReq) route.AppResp {
	var post model.Post
	jsonData, _ := json.Marshal(req.Body)
	json.Unmarshal(jsonData, &post)

	post.CreatorId = uint(req.Body["userId"].(float64))

	// Create the post in the database.
	err := u.Store.CreatePost(&post)
	if err != nil {
		return map[string]interface{}{
			"status": http.StatusNotFound,
			"error":  err.Error(),
		}
	}
	return map[string]interface{}{
		"status": http.StatusOK,
		"post":   post,
	}
}

// GetUsers retrieves all users from the database.
func (u *Post) GetPosts(ctx context.Context, req *route.AppReq) route.AppResp {
	var posts []model.Post
	err := u.Store.GetAllPosts(&posts)
	if err != nil {
		return map[string]interface{}{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		}
	}
	return map[string]interface{}{
		"status": http.StatusOK,
		"posts":  posts,
	}
}

// GetUser retrieves a user from the cache if present, else retrieves from the database.
func (ps *Post) GetPost(ctx context.Context, req *route.AppReq) route.AppResp {
	var post *model.Post
	id := req.Params["id"]
	postInterface, _ := ps.UserCache.Get("post.id=" + id)

	if postInterface == nil {
		post = &model.Post{}
		err := ps.Store.GetPost(post, "id="+id)
		if err != nil {
			return map[string]interface{}{
				"status": http.StatusInternalServerError,
				"error":  err.Error(),
			}
		}
		ps.UserCache.Set("post.id="+id, post, nil) // Set the post in the cache.
		return map[string]interface{}{
			"status": http.StatusOK,
			"post":   post,
		}
	} else {
		postInterface, _ := ps.UserCache.Get("post.id=" + id)
		postMap, ok := postInterface.(map[string]interface{})
		if !ok {
			return map[string]interface{}{
				"status": http.StatusInternalServerError,
				"error":  errors.New("postInterface is not mapStringInterface"),
			}
		}

		// Now you can unmarshal postMap into your user struct.
		userBytes, err := json.Marshal(postMap)
		if err != nil {
			return map[string]interface{}{
				"status": http.StatusInternalServerError,
				"error":  errors.New("json Marshal failed!!"),
			}
		}

		var post *model.Post
		err = json.Unmarshal(userBytes, &post)
		if err != nil {
			return map[string]interface{}{
				"status": http.StatusInternalServerError,
				"error":  errors.New("JSON unmarshal failed"),
			}
		}
		return map[string]interface{}{
			"status": http.StatusOK,
			"post":   post,
		}
	}
}

// deletePost removes a post from the database based on the given post ID.
func (u *Post) DeletePost(ctx context.Context, req *route.AppReq) route.AppResp {
	var post model.Post
	id := req.Params["id"]
	err := u.Store.GetPost(&post, "id="+id)
	if err != nil {
		return map[string]interface{}{
			"status": http.StatusNotFound,
			"error":  err.Error(),
		}
	}

	userId := req.Body["userId"]
	userRole := req.Body["role"]

	if userRole != "admin" && userId != post.CreatorId {
		return map[string]interface{}{
			"status":  http.StatusUnauthorized,
			"message": ("You are not allowed to Delete this post"),
		}
	}

	jsonData, _ := json.Marshal(req.Body)
	json.Unmarshal(jsonData, &post)

	err = u.Store.DeletePost(&post, id)
	if err != nil {
		return map[string]interface{}{
			"status": http.StatusNotFound,
			"error":  err.Error(),
		}
	}
	return map[string]interface{}{
		"status":  http.StatusOK,
		"message": "post with " + id + " is Deleted!",
	}
}

// UpdatePost updates a post in the database based on the given post ID.
func (ps *Post) UpdatePost(ctx context.Context, req *route.AppReq) route.AppResp {
	id := req.Params["id"]
	var post model.Post

	err := ps.Store.GetPost(&post, "id="+id)

	if err != nil {
		return map[string]interface{}{
			"status": http.StatusNotFound,
			"error":  err.Error(),
		}
	}
	userId := req.Body["userId"]
	userRole := req.Body["role"]

	if userRole != "admin" && userId != post.CreatorId {
		return map[string]interface{}{
			"status":  http.StatusUnauthorized,
			"message": ("You are not allowed to update this post"),
		}
	}

	jsonData, _ := json.Marshal(req.Body)
	json.Unmarshal(jsonData, &post)

	err = ps.Store.UpdatePost(&post)
	if err != nil {
		return map[string]interface{}{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		}
	}
	return map[string]interface{}{
		"status":  http.StatusOK,
		"message": "Post updated !!",
		"post":    post,
	}
}

// Add Comment add comment to  a post in the database based on the given post ID.
func (ps *Post) AddComment(ctx context.Context, req *route.AppReq) route.AppResp {
	id := req.Body["id"].(string)
	var post model.Post
	err := ps.Store.GetPost(&post, "id="+id)
	if err != nil {
		return map[string]interface{}{
			"status": http.StatusNotFound,
			"error":  err.Error(),
		}
	}
	var newComment model.Comment
	jsonData, _ := json.Marshal(req.Body["comment"])
	json.Unmarshal(jsonData, &newComment)
	post.Comments = append(post.Comments, newComment)

	err = ps.Store.UpdatePost(&post)
	if err != nil {
		return map[string]interface{}{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		}
	}
	return map[string]interface{}{
		"status":  http.StatusOK,
		"message": "comment added  !!",
		"comment": newComment,
		"post":    post,
	}
}
