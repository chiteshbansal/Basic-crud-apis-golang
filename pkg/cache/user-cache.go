package cache

import model "first-api/internal/models"

type UserCache interface {
	Set(key string, value *model.User)
	Get(key string) *model.User
}
