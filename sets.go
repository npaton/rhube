package rhube

import (
// "log"

// "strconv"
)

func (db *DB) Sadd(key string, members ...string) int {
	set := db.SetsMap[key]
	if set == nil {
		db.SetsMap[key] = make(map[string]bool)
		set = db.SetsMap[key]
	}

	addedCount := 0
	for member := range members {
		if set[members[member]] == false {
			addedCount++
			set[members[member]] = true
		}
	}

	return addedCount
}

func (db *DB) Srem(key string, members ...string) int {
	set := db.SetsMap[key]
	if set == nil {
		return 0
	}

	removedCount := 0
	for member := range members {
		if set[members[member]] == true {
			removedCount++
			delete(db.SetsMap[key], members[member])
		}
	}

	return removedCount
}

func (db *DB) Sismember(key string, member string) bool {
	set := db.SetsMap[key]
	if set == nil {
		return false
	}

	return set[member]
}

func (db *DB) Scard(key string) int {
	set := db.SetsMap[key]
	if set == nil {
		return 0
	}

	return len(set)
}

func (db *DB) Sinter(keys ...string) []string {
	var result []string
	shortest := db.SetsMap[keys[0]]
	shortestLength := len(shortest)
	for key := range keys {
		length := len(db.SetsMap[keys[key]])
		if length < shortestLength {
			shortest = db.SetsMap[keys[key]]
			shortestLength = length
		}
	}

	var skip bool
	for member, _ := range shortest {

		skip = false
		for key := range keys {
			if skip {
				continue
			}

			other := db.SetsMap[keys[key]]
			if other[member] == false {
				skip = true
			}
		}

		if !skip {
			result = append(result, member)
		}
	}

	return result
}

func (db *DB) Sinterstore(destination string, keys ...string) int {
	vals := db.Sinter(keys...)
	set := make(map[string]bool, 0)
	var count int
	for i, l := 0, len(vals); i < l; i++ {
		set[vals[i]] = true
		count++
	}
	db.SetsMap[destination] = set
	return count
}

func (db *DB) Smembers(key string) []string {
	set := db.SetsMap[key]
	if set == nil {
		return nil
	}

	members := make([]string, len(set))
	i := 0
	for k, _ := range set {
		members[i] = k
		i++
	}

	return members
}

func (db *DB) Smove(source, destination, member string) bool {
	setSource := db.SetsMap[source]
	if setSource == nil || setSource[member] == false {
		return false
	}

	if db.SetsMap[destination] == nil {
		db.SetsMap[destination] = make(map[string]bool)
	}

	db.SetsMap[destination][member] = true
	delete(db.SetsMap[source], member)

	return true
}

func (db *DB) Srandmember(key string) string {
	set := db.SetsMap[key]
	if set == nil {
		return ""
	}
	var result string
	for result, _ = range set {
		break
	}

	return result
}

func (db *DB) Spop(key string) string {
	member := db.Srandmember(key)
	db.Srem(key, member)
	return member
}

func (db *DB) Sunion(keys ...string) []string {
	var result []string
	seen := make(map[string]bool)

	for key := range keys {
		for val, _ := range db.SetsMap[keys[key]] {
			if seen[val] == false {
				result = append(result, val)
				seen[val] = true
			}
		}
	}

	return result
}

func (db *DB) Sunionstore(destination string, keys ...string) int {
	vals := db.Sunion(keys...)
	set := make(map[string]bool, 0)
	
	var count int
	for i, l := 0, len(vals); i < l; i++ {
		set[vals[i]] = true
		count++
	}
	
	db.SetsMap[destination] = set
	return count
}
