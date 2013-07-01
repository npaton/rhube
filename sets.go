package rhube

// Adds one or more members to a set
func (db *DB) Sadd(key string, members ...string) int {
	db.validateKeyType(key, "set")

	set := db.SetsMap[key]
	if set == nil {
		db.SetsMap[key] = make(map[string]bool)
		set = db.SetsMap[key]
	}

	addedCount := 0
	for _, member := range members {
		if set[member] == false {
			addedCount++
			set[member] = true
		}
	}

	return addedCount
}

func (db *DB) Srem(key string, members ...string) int {
	db.validateKeyType(key, "set")

	set := db.SetsMap[key]
	if set == nil {
		return 0
	}

	removedCount := 0
	for _, member := range members {
		if set[member] == true {
			removedCount++
			delete(db.SetsMap[key], member)
		}
	}

	return removedCount
}

func (db *DB) Sismember(key string, member string) bool {
	db.validateKeyType(key, "set")

	set := db.SetsMap[key]
	if set == nil {
		return false
	}

	return set[member]
}

func (db *DB) Scard(key string) int {
	db.validateKeyType(key, "set")

	set := db.SetsMap[key]
	if set == nil {
		return 0
	}

	return len(set)
}

func (db *DB) Sinter(keys ...string) []string {
	for _, key := range keys {
		db.validateKeyType(key, "set")	
	}

	var result []string
	shortest := db.SetsMap[keys[0]]
	shortestLength := len(shortest)
	for _, key := range keys {
		length := len(db.SetsMap[key])
		if length < shortestLength {
			shortest = db.SetsMap[key]
			shortestLength = length
		}
	}

	var skip bool
	for member, _ := range shortest {
		skip = false
		for _, key := range keys {
			if skip {
				continue
			}

			other := db.SetsMap[key]
			if other[member] == false {
				skip = true
				break
			}
		}

		if !skip {
			result = append(result, member)
		}
	}

	return result
}

func (db *DB) Sinterstore(destination string, keys ...string) int {
	db.validateKeyType(destination, "set")	
	for _, key := range keys {
		db.validateKeyType(key, "set")	
	}

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
	db.validateKeyType(key, "set")

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
	db.validateKeyType(source, "set")	
	db.validateKeyType(destination, "set")	
	db.validateKeyType(member, "set")

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
	db.validateKeyType(key, "set")

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
	db.validateKeyType(key, "set")

	member := db.Srandmember(key)
	db.Srem(key, member)
	return member
}

func (db *DB) Sunion(keys ...string) []string {
	for _, key := range keys {
		db.validateKeyType(key, "set")	
	}

	var result []string
	seen := make(map[string]bool)

	for _, key := range keys {
		for val, _ := range db.SetsMap[key] {
			if seen[val] == false {
				result = append(result, val)
				seen[val] = true
			}
		}
	}

	return result
}

func (db *DB) Sunionstore(destination string, keys ...string) int {
	db.validateKeyType(destination, "set")	
	for _, key := range keys {
		db.validateKeyType(key, "set")	
	}


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
