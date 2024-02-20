package entity

import "time"

type Storage interface {
	Set(key string, value string, expiration time.Duration) error
	Get(key string) (string, error)
	Update(key string, value string) error
}
