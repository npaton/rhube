package rhube

import (
	"testing"
)

func TestLsetLindex(t *testing.T) {
	db := NewDB()
	err := db.Lset("toto", 0, []byte("lala"))
	if err != nil {
		t.Fatal(err.Error())
	}

	if string(db.Lindex("toto", 0)) != "lala" {
		t.Fatal("should be lala", string(db.Lindex("toto", 0)))
	}

	err = db.Lset("toto", 42, []byte("hum"))
	if err != nil {
		t.Fatal(err.Error())
	}

	if string(db.Lindex("toto", 42)) != "hum" {
		t.Fatal("should be hum", string(db.Lindex("toto", 42)))
	}

	if string(db.Lindex("toto", 53)) != "" {
		t.Fatal("should be nil", string(db.Lindex("toto", 53)))
	}
}

func TestLinsert(t *testing.T) {
	db := NewDB()
	err := db.Lset("toto", 0, []byte("lala"))
	if err != nil {
		t.Fatal(err.Error())
	}

	err = db.Lset("toto", 1, []byte("tata"))
	if err != nil {
		t.Fatal(err.Error())
	}

	err = db.Lset("toto", 2, []byte("fafa"))
	if err != nil {
		t.Fatal(err.Error())
	}

	n, err := db.Linsert("toto", "AFTER", []byte("tata"), []byte("haha"))
	if n != 4 {
		t.Fatalf("n should equal 4: %d", n)
	}

	if err != nil {
		t.Fatalf(err.Error())
	}

	if string(db.Lindex("toto", 2)) != "haha" || string(db.Lindex("toto", 3)) != "fafa" {
		t.Fatal("should be haha and fafa:", string(db.Lindex("toto", 2)), string(db.Lindex("toto", 3)))
	}

	n, err = db.Linsert("toto", "BEFORE", []byte("tata"), []byte("zaza"))
	if n != 5 {
		t.Fatalf("n should equal 5: %d", n)
	}

	if err != nil {
		t.Fatalf(err.Error())
	}

	if string(db.Lindex("toto", 1)) != "zaza" || string(db.Lindex("toto", 2)) != "tata" || string(db.Lindex("toto", 3)) != "haha" {
		t.Fatal("should be haha and fafa:", string(db.Lindex("toto", 1)), string(db.Lindex("toto", 2)), string(db.Lindex("toto", 3)))
	}

}
