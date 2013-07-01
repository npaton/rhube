package rhube

import (
	"fmt"
	"math"
	"strconv"
	// "github.com/hotei/bits"
)

func (db *DB) Get(key string) []byte {
	return db.StringsMap[key]
}

func (db *DB) Set(key string, val []byte) bool {
	oldValue, exists := db.StringsMap[key]
	if exists && string(oldValue) == string(val) {
		return true
	}

	db.Del(key) // Set overrides anything
	db.StringsMap[key] = val
	return true
}

func (db *DB) Getset(key string, newValue []byte) ([]byte, error) {
	val := db.StringsMap[key]
	if string(val) == string(newValue) {
		return val, nil
	}

	db.Del(key) // Set overrides anything
	db.StringsMap[key] = newValue
	return val, nil
}

func (db *DB) Getbit(key string, offset int) int {
	db.validateKeyType(key, "string")
	val := db.StringsMap[key]
	if val == nil {
		return 0
	}
	pos := math.Ceil(float64(offset) / 8.0)
	if len(val) < int(pos) {
		return 0
	}

	offsetRem := math.Remainder(float64(val[int(pos-1)]), 8.0)
	return int(val[int(pos-1)]) & int(offsetRem)
}

func (db *DB) Setbit(key string, offset int, bit bool) int {
	db.validateKeyType(key, "string")
	original := db.Getbit(key, offset)
	val := db.StringsMap[key]

	pos := math.Ceil(float64(offset) / 8.0)
	if len(val) < int(pos) {
		for len(val) < int(pos) {
			val = append(val, byte(0x00))
		}
	}
	offsetRem := math.Remainder(float64(val[int(pos-1)]), 1.0)
	fmt.Println("pos", offset, pos, string(val), offsetRem, bit, val[int(pos-1)], float64(val[int(pos-1)]))
	v := val[int(pos-1)]
	if bit {
		v = byte(v | byte(offsetRem))
	} else {
		v = byte(v &^ byte(offsetRem))
	}
	val[int(pos-1)] = v

	db.StringsMap[key] = val
	return original
}

func (db *DB) Decrby(key string, decrement int) (int, error) {
	db.validateKeyType(key, "string")

	val := db.StringsMap[key]
	if val == nil {
		db.StringsMap[key] = []byte(strconv.Itoa(-decrement))
		return -decrement, nil
	}

	valInt, err := strconv.Atoi(string(val))
	if err != nil {
		return 0, err
	}
	valInt -= decrement
	db.StringsMap[key] = []byte(strconv.Itoa(valInt))
	return valInt, nil
}

func (db *DB) Decr(key string) (int, error) {
	return db.Decrby(key, 1)
}

func (db *DB) Incrby(key string, increment int) (int, error) {
	db.validateKeyType(key, "string")

	val := db.StringsMap[key]
	if val == nil {
		db.StringsMap[key] = []byte(strconv.Itoa(increment))
		return increment, nil
	}

	valInt, err := strconv.Atoi(string(val))
	if err != nil {
		return 0, err
	}
	valInt += increment
	db.StringsMap[key] = []byte(strconv.Itoa(valInt))
	return valInt, nil
}

func (db *DB) Incr(key string) (int, error) {
	return db.Incrby(key, 1)
}

func (db *DB) Incrbyfloat(key string, increment float64) (string, error) {
	db.validateKeyType(key, "string")

	val := db.StringsMap[key]
	if val == nil {
		str := strconv.FormatFloat(increment, 'f', -1, 64)
		db.StringsMap[key] = []byte(str)
		return str, nil
	}

	valFloat, err := strconv.ParseFloat(string(val), 64)
	if err != nil {
		return "", err
	}
	valFloat += increment
	str := strconv.FormatFloat(valFloat, 'f', -1, 64)
	db.StringsMap[key] = []byte(str)
	return str, nil
}

func (db *DB) Append(key string, val []byte) int {
	db.validateKeyType(key, "string")

	val = append(db.StringsMap[key], val...)
	db.StringsMap[key] = val
	return len(val)
}

func (db *DB) Strlen(key string) int {
	db.validateKeyType(key, "string")

	return len(db.StringsMap[key])
}

func (db *DB) Getrange(key string, start, stop int) []byte {
	db.validateKeyType(key, "string")

	val := db.StringsMap[key]
	length := len(val)
	if length == 0 {
		return []byte(val)
	}
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

func (db *DB) Setrange(key string, offset int, value []byte) int {
	db.validateKeyType(key, "string")

	val := db.StringsMap[key]
	length := len(val)

	if offset > length {
		padding := (offset - length) + len(value)
		for i := 0; i < padding; i++ {
			val = append(val, '\u0000')
		}
	}

	for i, l := 0, len(value); i < l; i++ {
		if cap(val) <= offset+i {
			val = append(val, value[i:]...)
			break
		} else {
			val[offset+i] = value[i]
		}
	}

	db.StringsMap[key] = val
	return len(db.StringsMap[key])
}

func (db *DB) Mget(keys ...string) [][]byte {
	result := make([][]byte, 0, len(keys))
	for i := range keys {
		val := db.StringsMap[keys[i]]
		result = append(result, val)
	}
	return result
}

func (db *DB) Mset(args ...string) {
	l := len(args)

	for i := 0; i < l/2; i++ {
		j := i * 2
		key, value := args[j], args[j+1]
		db.StringsMap[key] = []byte(value)
	}
}

func (db *DB) Msetnx(args ...string) bool {
	l := len(args)

	for i := 0; i < l/2; i++ {
		j := i * 2
		key := args[j]
		if db.StringsMap[key] != nil {
			return false
		}
	}

	for i := 0; i < l/2; i++ {
		j := i * 2
		key, value := args[j], args[j+1]
		db.StringsMap[key] = []byte(value)
	}

	return true
}
