package rhube

import (
	"testing"
)

func TestSaddSremSismember(t *testing.T) {
	db := NewDB()

	retVal := db.Sadd("alabama", "a", "b", "c")
	if retVal != 3 || db.Sismember("alabama", "random") || !db.Sismember("alabama", "a") || !db.Sismember("alabama", "b") || !db.Sismember("alabama", "c") {
		t.Fatal("Sadd: should add and find memebership", db.Sismember("alabama", "random"), !db.Sismember("alabama", "a"), !db.Sismember("alabama", "b"), !db.Sismember("alabama", "c"))
	}

	retVal = db.Sadd("alabama", "a", "b")
	if retVal != 0 {
		t.Fatal("Sadd: should add and ret proper changed count")
	}

	retVal = db.Srem("alabama", "a", "c", "random")
	if retVal != 2 || db.Sismember("alabama", "random") || db.Sismember("alabama", "a") || !db.Sismember("alabama", "b") || db.Sismember("alabama", "c") {
		t.Fatal("Sadd: should remove and find memebership", db.Sismember("alabama", "random"), !db.Sismember("alabama", "a"), !db.Sismember("alabama", "b"), !db.Sismember("alabama", "c"))
	}
}

func TestScard(t *testing.T) {
	db := NewDB()

	db.Sadd("alabama", "a", "b", "c")
	retVal := db.Scard("alabama")
	if retVal != 3 || db.Scard("kansas") != 0 {
		t.Fatal("Sadd: should return elements count 3 and 0:", retVal, db.Scard("kansas"))
	}
}

func TestSinter(t *testing.T) {
	db := NewDB()

	db.Sadd("alabama", "a", "b", "c")
	db.Sadd("florida", "c", "b", "c", "v")
	inter := db.Sinter("alabama", "florida")
	if inter[0] != "b" || inter[1] != "c" {
		t.Fatal("Sinter: should return b and c:", inter)
	}

	inter = db.Sinter("alabama", "kansas")
	if inter != nil {
		t.Fatal("Sinter: should return nothing:", inter)
	}
}

func TestSinterstore(t *testing.T) {
	db := NewDB()

	db.Sadd("alabama", "a", "b", "c")
	db.Sadd("florida", "c", "b", "c", "v")
	inter := db.Sinterstore("south", "alabama", "florida")
	if inter != 2 || !db.Sismember("south", "b") || !db.Sismember("south", "c") || db.Sismember("south", "a") {
		t.Fatal("Sinterstore: should store results:", inter)
	}

	inter = db.Sinterstore("neverland", "alabama", "kansas")
	if inter != 0 {
		t.Fatal("Sinter: should return nothing:", inter)
	}
}

func TestSunion(t *testing.T) {
	db := NewDB()

	db.Sadd("alabama", "a", "b", "c")
	db.Sadd("florida", "c", "b", "c", "v")
	union := db.Sunion("alabama", "florida")
	if union[0] != "a" || union[1] != "b" || union[2] != "c" || union[3] != "v" {
		t.Fatal("Sunion: should return a, b, c and v:", union)
	}

	union = db.Sunion("alabama", "kansas")
	if union[0] != "a" || union[1] != "b" || union[2] != "c" {
		t.Fatal("Sunion: should return nothing:", union)
	}
}

func TestSunionstore(t *testing.T) {
	db := NewDB()

	db.Sadd("alabama", "a", "b", "c")
	db.Sadd("florida", "c", "b", "c", "v")
	union := db.Sunionstore("south", "alabama", "florida")
	if union != 4 || !db.Sismember("south", "b") || !db.Sismember("south", "c") || !db.Sismember("south", "a") || !db.Sismember("south", "v") {
		t.Fatal("Sunionstore: should store results:", union)
	}

	union = db.Sunionstore("neverland", "joinville", "kansas")
	if union != 0 {
		t.Fatal("Sunion: should return nothing:", union)
	}
}

func TestSmembers(t *testing.T) {
	db := NewDB()

	db.Sadd("alabama", "a", "b", "c")
	inter := db.Smembers("alabama")
	if inter[0] != "a" || inter[1] != "b" || inter[2] != "c" {
		t.Fatal("Smembers: should return a, b and c:", inter)
	}

	inter = db.Smembers("kansas")
	if inter != nil {
		t.Fatal("Smembers: should return nothing:", inter)
	}
}

func TestSmove(t *testing.T) {
	db := NewDB()

	db.Sadd("alabama", "a", "b", "c")
	db.Sadd("mississippi", "d", "e", "f")
	retVal := db.Smove("alabama", "mississippi", "a")
	if retVal != true || !db.Sismember("mississippi", "a") || db.Sismember("alabama", "a") {
		t.Fatal("Smove: did not move", retVal)
	}

	db.Sadd("alabama", "a", "b", "c")
	retVal = db.Smove("alabama", "kansas", "x")
	if retVal != false || db.Sismember("kansas", "x") || db.Sismember("alabama", "x") {
		t.Fatal("Smove: should return nothing:", retVal, db.Sismember("kansas", "x"), db.Sismember("alabama", "x"))
	}
}

func TestSrandmember(t *testing.T) {
	db := NewDB()

	db.Sadd("alabama", "a", "b", "c")
	retVal := db.Srandmember("alabama")
	if retVal != "a" && retVal != "b" && retVal != "c" {
		t.Fatal("Srandmember: did not give a rand member inputed", retVal)
	}

	retVal = db.Srandmember("florida")
	if retVal != "" {
		t.Fatal("Smove: should return nothing:", retVal, db.Sismember("kansas", "x"), db.Sismember("alabama", "x"))
	}
}

func TestSpop(t *testing.T) {
	db := NewDB()

	db.Sadd("alabama", "a", "b", "c")
	retVal := db.Spop("alabama")
	if retVal != "a" && retVal != "b" && retVal != "c" {
		t.Fatal("Spop: did not give a rand member inputed", retVal)
	}

	match := 0
	for _, part := range []string{"a", "b", "c"} {
		if db.Sismember("alabama", part) {
			match++
		}
	}
	if match != 2 {
		t.Fatal("Spop: did not give remove member properly", retVal)
	}
}
