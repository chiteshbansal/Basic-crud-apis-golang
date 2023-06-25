package cache

import (
	"time"

	"github.com/go-redsync/redsync/v4"
)

type UserCache interface {
	Set(key string, value interface{}, exp *time.Duration) error
	Get(key string) (interface{}, error)
	GetMutex(mutexName string) *redsync.Mutex
}
