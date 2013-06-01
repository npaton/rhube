package rhube

import (
	// "log"
	"strconv"
)

func (db *DB) Get(key string) []byte {
	return db.KeyValuePairs[key]
}

func (db *DB) Set(key string, val []byte) bool {
	db.KeyValuePairs[key] = val
	return true
}

func (db *DB) Append(key string, val []byte) int {
	val = append(db.KeyValuePairs[key], val...)
	db.KeyValuePairs[key] = val
	return len(val)
}

func (db *DB) Decrby(key string, decrement int) (int, error) {
	val := db.KeyValuePairs[key]
	valInt, err := strconv.Atoi(string(val))
	if err != nil {
		return 0, err
	}
	valInt -= decrement
	db.KeyValuePairs[key] = []byte(strconv.Itoa(valInt))
	return valInt, nil
}

func (db *DB) Decr(key string) (int, error) {
	return db.Decrby(key, 1)
}

func (db *DB) Incrby(key string, increment int) (int, error) {
	val := db.KeyValuePairs[key]
	valInt, err := strconv.Atoi(string(val))
	if err != nil {
		return 0, err
	}
	valInt += increment
	db.KeyValuePairs[key] = []byte(strconv.Itoa(valInt))
	return valInt, nil
}

func (db *DB) Incr(key string) (int, error) {
	return db.Incrby(key, 1)
}

func (db *DB) Getrange(key string, start, stop int) []byte {
	val := db.KeyValuePairs[key]
	length := len(val)
	if stop < 0 {
		stop = length + (stop + 1)
	}
	if start < 0 {
		start = length + (start + 1)
	}
	if start > stop {
		return nil
	}
	if stop > length {
		stop = length
	}
	if start > length {
		start = length
	}
	return val[start:stop]
}
