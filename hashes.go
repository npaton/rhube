package rhube

import (
	"strconv"
)

func (db *DB) Hget(key, field string) string {
	db.validateKeyType(key, "hash")
	
	hash := db.HashesMap[key]
	if hash == nil {
		return ""
	}

	if hash[field] == "" {
		return ""
	}
	return hash[field]
}

func (db *DB) Hset(key, field, value string) bool {
	db.validateKeyType(key, "hash")
	
	hash := db.HashesMap[key]
	newField := false

	if hash == nil {
		hash = make(map[string]string)
		newField = true
	}

	if hash[field] == "" {
		newField = true
	}

	hash[field] = value
	db.HashesMap[key] = hash

	return newField
}

func (db *DB) Hsetnx(key, field, value string) bool {
	db.validateKeyType(key, "hash")
	
	hash := db.HashesMap[key]

	if hash == nil {
		hash = make(map[string]string)
	}

	if hash[field] != "" {
		return false
	}

	hash[field] = value
	db.HashesMap[key] = hash

	return true
}

func (db *DB) Hmset(key string, args ...string) bool {
	db.validateKeyType(key, "hash")
	
	l := len(args)

	for i := 0; i < l/2; i++ {
		j := i * 2
		field, value := args[j], args[j+1]
		hash := db.HashesMap[key]
		if hash == nil {
			hash = make(map[string]string)
		}
		hash[field] = value
		db.HashesMap[key] = hash
	}
	return true
}

func (db *DB) Hmget(key string, fields ...string) []string {
	db.validateKeyType(key, "hash")
	
	result := make([]string, 0, len(fields))
	hash := db.HashesMap[key]
	if hash == nil {
		return nil
	}

	for _, field := range fields {
		result = append(result, hash[field])
	}

	return result
}

func (db *DB) Hexist(key string, field string) bool {
	db.validateKeyType(key, "hash")
	
	hash := db.HashesMap[key]
	if hash == nil {
		return false
	}

	return hash[field] != ""
}

func (db *DB) Hdel(key string, fields ...string) int {
	db.validateKeyType(key, "hash")
	
	hash := db.HashesMap[key]
	if hash == nil {
		return 0
	}
	fieldsChanged := 0
	for i, l := 0, len(fields); i < l; i++ {
		if hash[fields[i]] != "" {
			delete(db.HashesMap[key], fields[i])
			fieldsChanged++
		}
	}
	if len(db.HashesMap[key]) == 0 {
		delete(db.HashesMap, key)
		db.cancelExpireKey(key)
	}
	return fieldsChanged
}

func (db *DB) Hgetall(key string) map[string]string {
	db.validateKeyType(key, "hash")
	
	hash := db.HashesMap[key]
	if hash == nil {
		return nil
	}

	return hash
}

func (db *DB) Hincrby(key, field string, increment int) (int, error) {
	db.validateKeyType(key, "hash")
	
	hash := db.HashesMap[key]
	if hash == nil {
		hash = make(map[string]string)
	}

	if hash[field] == "" {
		hash[field] = "0"
	}

	intVal, err := strconv.Atoi(string(hash[field]))
	if err != nil {
		return 0, err
	}

	intVal += increment
	db.HashesMap[key][field] = strconv.Itoa(intVal)
	return intVal, nil
}

func (db *DB) Hincrbyfloat(key, field string, increment float64) (string, error) {
	db.validateKeyType(key, "hash")
	
	hash := db.HashesMap[key]
	if hash == nil {
		hash = make(map[string]string)
	}

	if hash[field] == "" {
		hash[field] = "0"
	}

	intVal, err := strconv.ParseFloat(string(hash[field]), 64)
	if err != nil {
		return "", err
	}

	intVal += increment
	db.HashesMap[key][field] = strconv.FormatFloat(intVal, 'f', -1, 64)
	return db.HashesMap[key][field], nil
}

func (db *DB) Hkeys(key string) []string {
	db.validateKeyType(key, "hash")
	
	hash := db.HashesMap[key]
	if hash == nil {
		return nil
	}

	var result []string
	for key, _ := range hash {
		if key != "" {
			result = append(result, key)
		}
	}
	return result
}

func (db *DB) Hlen(key string) int {
	db.validateKeyType(key, "hash")
	
	hash := db.HashesMap[key]
	if hash == nil {
		return 0
	}

	return len(hash)
}

func (db *DB) Hvals(key string) []string {
	db.validateKeyType(key, "hash")
	
	hash := db.HashesMap[key]
	if hash == nil {
		return nil
	}

	var result []string
	for key, val := range hash {
		if key != "" && val != "" {
			result = append(result, val)
		}
	}
	return result
}
