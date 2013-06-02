package rhube

import (
	//	"log"
	"strconv"
)

func (db *DB) Get(key string) []byte {
	return db.StringsMap[key]
}

func (db *DB) Set(key string, val []byte) bool {
	db.StringsMap[key] = val
	return true
}

func (db *DB) Getset(key string, newValue []byte) ([]byte, error) {
	val := db.StringsMap[key]
	db.StringsMap[key] = newValue
	return val, nil
}

func (db *DB) Decrby(key string, decrement int) (int, error) {
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
	val = append(db.StringsMap[key], val...)
	db.StringsMap[key] = val
	return len(val)
}

func (db *DB) Strlen(key string) int {
	return len(db.StringsMap[key])
}

func (db *DB) Getrange(key string, start, stop int) []byte {
	val := db.StringsMap[key]
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

func (db *DB) Setrange(key string, offset int, value []byte) int {
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
		// log.Println(i, keys[i], val)
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

// type intValue int

// func (i *intValue) MarshalRhube() ([]byte, error) {
// 	return []byte(strconv.Itoa(int(*i))), nil
// }

// func (i *intValue) UnmarshalRhube(val []byte) error {
// 	a, err := strconv.Atoi(string(val))
// 	*i = intValue(a)
// 	return err
// }

// func TestIntValue(t *testing.T) {
// 	valInt := intValue(3)
// 	val, _ := valInt.MarshalRhube()
// 	intVal := intValue(10)
// 	intVal.UnmarshalRhube(val)
// 	if int(intVal) != 3 {
// 		t.Fatalf("should be 3, no? %+v", intVal)
// 	}
// }
