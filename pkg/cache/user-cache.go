package cache

import model "first-api/api/Models"

type UserCache interface {
	Set(key string, value *model.User)
	Get(key string) *model.User
}
