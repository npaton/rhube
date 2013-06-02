package rhube

import (
	// "log"
	"testing"
)

func TestHsetHget(t *testing.T) {
	db := NewDB()

	db.Hset("key1", "field1", "abc")
	if db.Hget("key1", "field1") != "abc" {
		t.Fatal("Hset/Hget:", "abc", db.Hget("key1", "field1"))
	}

	db.Hset("key1", "field2", "432")
	if db.Hget("key1", "field2") != "432" {
		t.Fatal("Hset/Hget:", "432", db.Hget("key1", "field2"))
	}

	db.Hset("key1", "field2", "876")
	if db.Hget("key1", "field2") != "876" {
		t.Fatal("Hset/Hget:", "876", db.Hget("key1", "field2"))
	}

	if db.Hget("key2", "field1") != "" {
		t.Fatal("Hset/Hget:", "", db.Hget("key2", "field1"))
	}
}

func TestHsetnx(t *testing.T) {
	db := NewDB()

	db.Hset("key1", "field1", "")

	retVal := db.Hsetnx("key1", "field1", "abc")
	if retVal != true || db.Hget("key1", "field1") != "abc" {
		t.Fatal("Hsetnx:", "abc", db.Hget("key1", "field1"))
	}

	db.Hsetnx("key1", "field2", "432")
	if retVal != true || db.Hget("key1", "field2") != "432" {
		t.Fatal("Hsetnx:", "432", db.Hget("key1", "field2"))
	}

	retVal = db.Hsetnx("key1", "field1", "876")
	if retVal != false || db.Hget("key1", "field2") != "432" {
		t.Fatal("Hsetnx:", "432", db.Hget("key1", "field2"))
	}
}

func TestHmsetHmget(t *testing.T) {
	db := NewDB()

	db.Hmset("key1", "field1", "abc", "field2", "543")
	val := db.Hmget("key1", "field1", "field2")
	if val[0] != "abc" || val[1] != "543" {
		t.Fatal("Hmset/Hmget:", "abc, 543", val)
	}

	db.Hmset("key1", "field3", "abc", "field4", "543")
	val = db.Hmget("key1", "field1", "field3")
	if val[0] != "abc" || val[1] != "abc" {
		t.Fatal("Hmset/Hmget:", "abc, abc", val)
	}

	db.Hmset("key1", "field1", "098", "field2", "890")
	val = db.Hmget("key1", "field1", "field2")
	if val[0] != "098" || val[1] != "890" {
		t.Fatal("Hmset/Hmget:", "098, 890", val)
	}

}

func TestHexist(t *testing.T) {
	db := NewDB()

	db.Hset("key1", "field1", "abc")
	if !db.Hexist("key1", "field1") {
		t.Fatal("Hexist:", "should exist")
	}

	if db.Hexist("key1", "field2") {
		t.Fatal("Hexist:", "should exist")
	}

}

func TestHdel(t *testing.T) {
	db := NewDB()

	db.Hset("key1", "field1", "abc")
	db.Hset("key1", "field2", "abc")
	db.Hset("key1", "field3", "abc")

	db.Hdel("key1", "field1", "field2")

	if db.Hexist("key1", "field1") || db.Hexist("key1", "field2") {
		t.Fatal("Hexist:", "should not exist")
	}

	if !db.Hexist("key1", "field3") {
		t.Fatal("Hexist:", "should exist", db.Hget("key1", "field3"))
	}
}

func TestHgetall(t *testing.T) {
	db := NewDB()

	db.Hset("key1", "field1", "abc")
	db.Hset("key1", "field2", "abc")
	db.Hset("key1", "field3", "abc")

	all := db.Hgetall("key1")

	if all["field1"] != "abc" || all["field2"] != "abc" || all["field3"] != "abc" {
		t.Fatal("Hgetall:", "should get all fields")
	}

	all = db.Hgetall("key2")
	if len(all) > 0 || all["field1"] != "" {
		t.Fatal("Hgetall:", "should be empty")
	}
}

func TestHincrby(t *testing.T) {
	db := NewDB()

	db.Hset("key1", "field1", "2")

	total, err := db.Hincrby("key1", "field1", 5)
	if total != 7 || err != nil || db.Hget("key1", "field1") != "7" {
		t.Fatal("key1 should be 7 ", total, db.Hget("key1", "field1"))
	}

	db.Hset("key1", "field1", "-207")

	total, err = db.Hincrby("key1", "field1", 5)
	if total != -202 || err != nil || db.Hget("key1", "field1") != "-202" {
		t.Fatal("key1 should be 202", total, db.Hget("key1", "field1"))
	}
}

func TestHincrbyfloat(t *testing.T) {
	db := NewDB()

	db.Hset("key1", "field1", "2.0")

	total, err := db.Hincrbyfloat("key1", "field1", 5.2)
	if total != "7.2" || err != nil || db.Hget("key1", "field1") != "7.2" {
		t.Fatal("key1 should be 7.2", total, db.Hget("key1", "field1"))

	}

	db.Hset("key1", "field1", "-207.7")

	total, err = db.Hincrbyfloat("key1", "field1", 5.1)
	if total != "-202.6" || err != nil || db.Hget("key1", "field1") != "-202.6" {
		t.Fatal("key1 should be 202.6", total, db.Hget("key1", "field1"))
	}
}

func TestHkeys(t *testing.T) {
	db := NewDB()

	db.Hset("key1", "field1", "abc")
	db.Hset("key1", "field2", "abc")
	db.Hset("key1", "field3", "abc")

	keys := db.Hkeys("key1")

	if keys[0] != "field1" || keys[1] != "field2" || keys[2] != "field3" {
		t.Fatalf("Hkeys:", "should get all fields %#v", keys[0], keys)
	}

	keys = db.Hkeys("key2")

	if keys != nil {
		t.Fatalf("Hkeys:", "should be nil")
	}
}

func TestHvals(t *testing.T) {
	db := NewDB()

	db.Hset("key1", "field1", "abc")
	db.Hset("key1", "field2", "def")
	db.Hset("key1", "field3", "ghi")

	keys := db.Hvals("key1")

	if keys[0] != "abc" || keys[1] != "def" || keys[2] != "ghi" {
		t.Fatalf("Hvals:", "should get all fields %#v", keys[0], keys)
	}

	keys = db.Hvals("key2")

	if keys != nil {
		t.Fatalf("Hvals:", "should be nil")
	}
}

func TestHlen(t *testing.T) {
	db := NewDB()

	db.Hset("key1", "field1", "abc")
	db.Hset("key1", "field2", "abc")
	db.Hset("key1", "field3", "abc")

	length := db.Hlen("key1")

	if length != 3 {
		t.Fatalf("Hkeys:", "should get fields count 3:", length)
	}

	length = db.Hlen("key2")

	if length != 0 {
		t.Fatalf("Hkeys:", "should be 0", length)
	}
}
