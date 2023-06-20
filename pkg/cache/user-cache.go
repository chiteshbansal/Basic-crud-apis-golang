package cache

import (
	"time"
)

type UserCache interface {
	Set(key string, value interface{}, exp *time.Duration) error
	Get(key string) (interface{}, error)
}
