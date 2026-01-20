package boltdb

import (
	"time"

	bolt "go.etcd.io/bbolt"
)

func Open(path string) (*bolt.DB, error) {
	return bolt.Open(path, 0o600, &bolt.Options{Timeout: 1 * time.Second})
}
