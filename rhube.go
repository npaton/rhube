package rhube

import (
	// "log"
)

type DB struct {
	KeyValuePairs map[string][]byte
}

func NewDB() *DB {
	return &DB{KeyValuePairs: make(map[string][]byte)}
}

