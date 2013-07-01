package rhube

import (
	"testing"
)

func TestDel(t *testing.T) {
	db := NewDB()
	db.Set("ho", []byte("12"))
	db.Hset("ha", "hi", "hu")
	db.Lset("he", 0, []byte("hehe"))
	db.Sadd("hou", "ouh")

	if string(db.Get("ho")) != "12" || db.Hget("ha", "hi") != "hu" || string(db.Lindex("he", 0)) != "hehe" || !db.Sismember("hou", "ouh") {
		t.Fatalf("Failed")
	}

	i := db.Del("ho")
	if i != 1 || string(db.Get("ho")) != "" || db.Hget("ha", "hi") != "hu" || string(db.Lindex("he", 0)) != "hehe" || !db.Sismember("hou", "ouh") {
		t.Fatalf("Failed")
	}

	i = db.Del("ha")
	if i != 1 || string(db.Get("ho")) != "" || db.Hget("ha", "hi") != "" || string(db.Lindex("he", 0)) != "hehe" || !db.Sismember("hou", "ouh") {
		t.Fatalf("Failed")
	}

	i = db.Del("he", "hou")
	if i != 2 || string(db.Get("ho")) != "" || db.Hget("ha", "hi") != "" || string(db.Lindex("he", 0)) != "" || db.Sismember("hou", "ouh") {
		t.Fatalf("Failed")
	}
}

func TestKeys(t *testing.T) {
	db := NewDB()
	db.Set("ho", []byte("12"))
	db.Hset("ha", "hi", "hu")
	db.Lset("he", 0, []byte("hehe"))
	db.Sadd("hou", "ouh")

	if string(db.Get("ho")) != "12" || db.Hget("ha", "hi") != "hu" || string(db.Lindex("he", 0)) != "hehe" || !db.Sismember("hou", "ouh") {
		t.Fatalf("Failed")
	}

	keys := db.Keys("h[ou]")
	if len(keys) != 1 || keys[0] != "ho" {
		t.Fatalf("Failed", len(keys), keys)
	}

	keys = db.Keys("ho*")
	if len(keys) != 2 || keys[0] != "ho" || keys[1] != "hou" {
		t.Fatalf("Failed", len(keys), keys)
	}

	keys = db.Keys("h*")
	if len(keys) != 4 {
		t.Fatalf("Failed", len(keys), keys)
	}

	keys = db.Keys("h?")
	if len(keys) != 3 {
		t.Fatalf("Failed", len(keys), keys)
	}

	keys = db.Keys("*e")
	if len(keys) != 1 || keys[0] != "he" {
		t.Fatalf("Failed", len(keys), keys)
	}

	keys = db.Keys("*o*")
	if len(keys) != 2 || keys[0] != "ho" || keys[1] != "hou" {
		t.Fatalf("Failed", len(keys), keys)
	}
}

func TestRename(t *testing.T) {
	db := NewDB()
	db.Set("ho", []byte("12"))
	db.Rename("ho", "ha")

	if string(db.Get("ho")) != "" || string(db.Get("ha")) != "12" {
		t.Fatalf("Failed", string(db.Get("ho")), string(db.Get("ha")))
	}

	err := db.Rename("ho", "ha")

	if err == nil {
		t.Fatalf("Error should not be nil")
	}

	if string(db.Get("ho")) != "" || string(db.Get("ha")) != "12" {
		t.Fatalf("Failed", string(db.Get("ho")), string(db.Get("ha")))
	}

	err = db.Rename("ha", "ha")

	if err == nil {
		t.Fatalf("Error should not be nil")
	}
}

func TestRenamenx(t *testing.T) {
	db := NewDB()
	db.Set("ho", []byte("12"))
	db.Renamenx("ho", "ha")

	if string(db.Get("ho")) != "" || string(db.Get("ha")) != "12" {
		t.Fatalf("Failed", string(db.Get("ho")), string(db.Get("ha")))
	}

	i, err := db.Renamenx("ho", "ha")

	if i != 0 || err == nil {
		t.Fatalf("Error should not be nil")
	}

	if string(db.Get("ho")) != "" || string(db.Get("ha")) != "12" {
		t.Fatalf("Failed", string(db.Get("ho")), string(db.Get("ha")))
	}

	i, err = db.Renamenx("ha", "ha")

	if i != 0 || err == nil {
		t.Fatalf("Error should not be nil")
	}

	db.Set("hi", []byte("42"))
	i, err = db.Renamenx("ha", "hi")

	if i != 0 || err != nil {
		t.Fatalf("Error should not be nil")
	}

}

func TestExists(t *testing.T) {
	db := NewDB()
	db.Set("ho", []byte("12"))
	db.Hset("ha", "hi", "hu")
	db.Lset("he", 0, []byte("hehe"))
	db.Sadd("hou", "ouh")

	if db.Exists("ho") != true {
		t.Fatalf("Failed")
	}

	if db.Exists("ha") != true {
		t.Fatalf("Failed")
	}

	if db.Exists("he") != true {
		t.Fatalf("Failed")
	}

	if db.Exists("hou") != true {
		t.Fatalf("Failed")
	}

	if db.Exists("ouh") != false {
		t.Fatalf("Failed")
	}
}

func TestRandomkey(t *testing.T) {
	db := NewDB()
	db.Set("ho", []byte("12"))
	db.Hset("ha", "hi", "hu")
	db.Lset("he", 0, []byte("hehe"))
	db.Sadd("hou", "ouh")
	counts := make(map[string]int)

	for i := 1000; i > 0; i-- {
		k := db.Randomkey()
		if !(k == "ho" || k == "ha" || k == "he" || k == "hou") {
			t.Fatalf("Actual random key!", k)
		}
		counts[k]++
	}

	for k, count := range counts {
		if count == 0 {
			t.Fatalf("key not being picked: %s", k)
		}
	}
}
