package rhube

import (
	"fmt"
)

func (db *DB) Lset(key string, index int, value []byte) error {
	db.validateKeyType(key, "list")
	
	if len(db.ListsMap[key]) < index+1 {
		t := make([][]byte, index+1, index+1)
		copy(t, db.ListsMap[key])
		db.ListsMap[key] = t
	}

	db.ListsMap[key][index] = value

	return nil
}

func (db *DB) Lindex(key string, index int) []byte {
	db.validateKeyType(key, "list")
	
	if len(db.ListsMap[key]) < index+1 {
		return nil
	}
	
	return db.ListsMap[key][index]
}

func (db *DB) Linsert(key string, beforeAfter string, pivot []byte, value []byte) (int, error) {
	db.validateKeyType(key, "list")

	if (beforeAfter != "BEFORE" && beforeAfter != "AFTER") {
		return 0, fmt.Errorf("Linsert: can only insert key be 'BEFORE' or 'AFTER' pivot, not: '%s'", beforeAfter)	
	}
	
	p := string(pivot)
	found := false 
	for i, item := range db.ListsMap[key] {
		if string(item) == p {
			insertionIndex := i
			if beforeAfter == "AFTER" {
				insertionIndex = i + 1
			}
			db.ListsMap[key] = append(db.ListsMap[key][:insertionIndex], append([][]byte{value}, db.ListsMap[key][insertionIndex:]...)...)
			found = true
			break
		}
	}
	
	if found {
		return len(db.ListsMap[key]), nil
	} else {
		return 0, fmt.Errorf("Linsert: pivot not found in list: '%s'", pivot)
	}
}
