package rhube

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

func (db *DB) Del(keys ...string) int {
	count := 0
	for _, key := range keys {
		if _, ok := db.StringsMap[key]; ok {
			delete(db.StringsMap, key)
			db.cancelExpireKey(key)
			count++
			continue
		}

		if _, ok := db.HashesMap[key]; ok {
			delete(db.HashesMap, key)
			db.cancelExpireKey(key)
			count++
			continue
		}

		if _, ok := db.SetsMap[key]; ok {
			delete(db.SetsMap, key)
			db.cancelExpireKey(key)
			count++
			continue
		}

		if _, ok := db.ListsMap[key]; ok {
			delete(db.ListsMap, key)
			db.cancelExpireKey(key)
			count++
			continue
		}
	}

	return count
}

func (db *DB) Keys(pattern string) []string {
	result := make([]string, 0)
	pattern = strings.Replace(pattern, "?", ".", -1)
	pattern = strings.Replace(pattern, "*", ".*", -1)
	pattern = "^" + pattern + "$"

	for key, _ := range db.StringsMap {
		matchKey(pattern, key, &result)
	}

	for key, _ := range db.HashesMap {
		matchKey(pattern, key, &result)
	}

	for key, _ := range db.SetsMap {
		matchKey(pattern, key, &result)
	}

	for key, _ := range db.ListsMap {
		matchKey(pattern, key, &result)
	}

	return result
}

func matchKey(pattern, key string, result *[]string) {
	match, err := regexp.Match(pattern, []byte(key))
	if err == nil && match {
		*result = append(*result, key)
	}
}

// TODO: Support for expires!
func (db *DB) Renamenx(key, newKey string) (int, error) {
	if key == newKey {
		return 0, fmt.Errorf("Rename to same key name: %s", key)
	}

	if _, ok := db.StringsMap[key]; ok {
		if _, ok = db.StringsMap[newKey]; !ok {
			db.StringsMap[newKey] = db.StringsMap[key]
			delete(db.StringsMap, key)
			db.renameExpire(key, newKey)
			return 1, nil
		} else {
			return 0, nil
		}
	}

	if _, ok := db.HashesMap[key]; ok {
		if _, ok = db.HashesMap[newKey]; !ok {
			db.HashesMap[newKey] = db.HashesMap[key]
			delete(db.HashesMap, key)
			db.renameExpire(key, newKey)
			return 1, nil
		} else {
			return 0, nil
		}
	}

	if _, ok := db.SetsMap[key]; ok {
		if _, ok = db.SetsMap[newKey]; !ok {
			db.SetsMap[newKey] = db.SetsMap[key]
			delete(db.SetsMap, key)
			db.renameExpire(key, newKey)
			return 1, nil
		} else {
			return 0, nil
		}
	}

	if _, ok := db.ListsMap[key]; ok {
		if _, ok = db.ListsMap[newKey]; !ok {
			db.ListsMap[newKey] = db.ListsMap[key]
			delete(db.ListsMap, key)
			db.renameExpire(key, newKey)
			return 1, nil
		} else {
			return 0, nil
		}
	}

	return 0, fmt.Errorf("Rename: '%s' not found", key)
}

// TODO: Support for expires!
func (db *DB) Rename(key, newKey string) error {
	if key == newKey {
		return fmt.Errorf("Rename to same key name: %s", key)
	}

	db.renameExpire(key, newKey)

	if _, ok := db.StringsMap[key]; ok {
		db.StringsMap[newKey] = db.StringsMap[key]
		delete(db.StringsMap, key)
		return nil
	}

	if _, ok := db.HashesMap[key]; ok {
		db.HashesMap[newKey] = db.HashesMap[key]
		delete(db.HashesMap, key)
		return nil
	}

	if _, ok := db.SetsMap[key]; ok {
		db.SetsMap[newKey] = db.SetsMap[key]
		delete(db.SetsMap, key)
		return nil
	}

	if _, ok := db.ListsMap[key]; ok {
		db.ListsMap[newKey] = db.ListsMap[key]
		delete(db.ListsMap, key)
		return nil
	}

	return fmt.Errorf("Source Key not found: %s", key)
}

func (db *DB) Exists(key string) bool {

	if _, ok := db.StringsMap[key]; ok {
		return true
	}

	if _, ok := db.HashesMap[key]; ok {
		return true
	}

	if _, ok := db.SetsMap[key]; ok {
		return true
	}

	if _, ok := db.ListsMap[key]; ok {
		return true
	}

	return false
}

// Aweful implementation...
func (db *DB) Randomkey() string {
	maps := make([]interface{}, 4)
	maps = append(maps, db.StringsMap)
	maps = append(maps, db.HashesMap)
	maps = append(maps, db.SetsMap)
	maps = append(maps, db.ListsMap)

	if len(db.StringsMap) == 0 && len(db.HashesMap) == 0 && len(db.SetsMap) == 0 && len(db.ListsMap) == 0 {
		return ""
	}

	result := ""
	for loops := 100; loops != 0 || result == ""; loops-- {
		collN := rand.Intn(len(maps) / 2)

		switch collN {
		case 0:
			itemN := rand.Intn(len(db.StringsMap))
			count := 0
			for key, _ := range db.StringsMap {
				if count == itemN {
					return key
				}
				count++
			}
		case 1:
			itemN := rand.Intn(len(db.HashesMap))
			count := 0
			for key, _ := range db.HashesMap {
				if count == itemN {
					return key
				}
				count++
			}
		case 2:
			itemN := rand.Intn(len(db.SetsMap))
			count := 0
			for key, _ := range db.SetsMap {
				if count == itemN {
					return key
				}
				count++
			}
		case 3:
			itemN := rand.Intn(len(db.ListsMap))
			count := 0
			for key, _ := range db.ListsMap {
				if count == itemN {
					return key
				}
				count++
			}
		default:
			fmt.Println("DEFAUTL, should not be here")
		}
	}

	return result
}











// Tests needed...





func (db *DB) Persist(key string) bool {
	return db.cancelExpireKey(key)
}

func (db *DB) Expire(key string, seconds int) bool {
	return db.Pexpire(key, seconds*1000)
}

func (db *DB) Pexpire(key string, milliseconds int) bool {
	if !db.Exists(key) {
		return false
	}

	db.expireKeyIn(key, milliseconds)
	return true
}

func (db *DB) Expireat(key string, t time.Time) bool {
	if !db.Exists(key) {
		return false
	}

	db.expireKeyAt(key, t)
	return true
}

func (db *DB) Pexpireat(key string, t time.Time) bool {
	if !db.Exists(key) {
		return false
	}

	db.expireKeyAt(key, t)
	return true
}

func (db *DB) TTL(key string) int {
	if !db.Exists(key) {
		return -2
	}

	e, exist := db.ExpiresMap[key]
	if !exist {
		return -1
	}

	return e.TTL()
}

func (db *DB) PTTL(key string) int {
	if !db.Exists(key) {
		return -2
	}

	e, exist := db.ExpiresMap[key]
	if !exist {
		return -1
	}

	return e.PTTL()
}

func (db *DB) Type(key string) string {

	if _, ok := db.StringsMap[key]; ok {
		return "string"
	}

	if _, ok := db.HashesMap[key]; ok {
		return "hash"
	}

	if _, ok := db.SetsMap[key]; ok {
		return "set"
	}

	if _, ok := db.ZsetsMap[key]; ok {
		return "zset"
	}

	if _, ok := db.ListsMap[key]; ok {
		return "list"
	}

	return "none"
}

func (db *DB) validateKeyType(key string, expected string) error {

	t := db.Type(key)
	if t != expected && t != "none" {
		panic(fmt.Errorf("Operation against a key holding the wrong kind of value (current type: %s)", t))
	}

	return nil
}
