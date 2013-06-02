package rhube

import (
// "log"
// "io"
)

type DB struct {
	StringsMap map[string][]byte
	HashesMap  map[string]map[string]string
	SetsMap    map[string]map[string]bool
}

func NewDB() *DB {
	return &DB{
		StringsMap: make(map[string][]byte),
		HashesMap:  make(map[string]map[string]string),
		SetsMap:    make(map[string]map[string]bool),
	}
}

type Value interface {
	MarshalRhube() ([]byte, error)
	UnmarshalRhube([]byte) error
}
