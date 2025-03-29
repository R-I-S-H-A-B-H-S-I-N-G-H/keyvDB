package db

import (
	"errors"
	"time"
)

type DbVal struct {
	Val    string
	Expiry int64
}

type Db struct {
	KeyVal map[string]*DbVal
}

type Store struct {
	db *Db
}

func NewDb() *Db {
	return &Db{KeyVal: make(map[string]*DbVal)}
}

func NewStore(db *Db) *Store {
	return &Store{db: db}
}

func (s *Store) Get(key string) (string, error) {
	val, exists := s.db.KeyVal[key]
	if !exists {
		return "", errors.New("key not found")
	}

	// Check if the key has expired
	if isExpired(val.Expiry) {
		delete(s.db.KeyVal, key) // Remove expired key
		return "", errors.New("key expired")
	}

	return val.Val, nil
}

// Set stores a key-value pair with an optional expiration time (in seconds).
func (s *Store) Set(key string, value string, expiryInSec int64) {
	expiryTimeInNano := generateExirationTime(convertSecToNano(expiryInSec))

	s.db.KeyVal[key] = &DbVal{
		Val:    value,
		Expiry: expiryTimeInNano,
	}
}

func isExpired(futurTimeInNano int64) bool {
	if futurTimeInNano < 0 {
		return false
	}

	currentTimeNano := getCurrentTimeNano()

	return futurTimeInNano < currentTimeNano
}

func generateExirationTime(expiryInNano int64) int64 {
	if expiryInNano < 0 {
		return -1
	}
	return getCurrentTimeNano() + expiryInNano
}

func getCurrentTimeNano() int64 {
	return time.Now().UnixNano()
}

func convertSecToNano(sec int64) int64 {
	return sec * time.Second.Nanoseconds()
}
