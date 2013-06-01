package rhube

import (
	"testing"
	// "time"
	"strconv"
)

func TestGetSet(t *testing.T) {
	db := NewDB()
	db.Set("ho", []byte("12"))
	val := db.Get("ho")
	if string(val) != "12" {
		t.Fatalf("Wierd!")
	}
}

func TestAppend(t *testing.T) {
	db := NewDB()
	db.Set("ho", []byte("12"))
	length := db.Append("ho", []byte("12"))
	val := db.Get("ho")
	if string(val) != "1212" || length != 4 {
		t.Fatalf("Wierd!")
	}
}

func TestDecr(t *testing.T) {
	db := NewDB()
	db.Set("ho", []byte("12"))

	length, err := db.Decr("ho")
	val := db.Get("ho")
	valInt, err := strconv.Atoi(string(val))
	if err != nil {
		t.Fatalf(err.Error())
	}
	if length != valInt {
		t.Fatalf("before after not equal")
	}

	if valInt != 11 {
		t.Fatalf("after not equal to 11")
	}

	length, err = db.Decr("ho")
	val = db.Get("ho")
	valInt, err = strconv.Atoi(string(val))
	if err != nil {
		t.Fatalf(err.Error())
	}
	if length != valInt {
		t.Fatalf("before after not equal")
	}

	if valInt != 10 {
		t.Fatalf("after not equal to 11")
	}

}

func TestDecrby(t *testing.T) {
	db := NewDB()
	db.Set("ho", []byte("12"))

	length, err := db.Decrby("ho", 2)
	val := db.Get("ho")
	valInt, err := strconv.Atoi(string(val))
	if err != nil {
		t.Fatalf(err.Error())
	}
	if length != valInt {
		t.Fatalf("before after not equal")
	}

	if valInt != 10 {
		t.Fatalf("after not equal to 10")
	}

	length, err = db.Decrby("ho", 5)
	val = db.Get("ho")
	valInt, err = strconv.Atoi(string(val))
	if err != nil {
		t.Fatalf(err.Error())
	}
	if length != valInt {
		t.Fatalf("before after not equal")
	}

	if valInt != 5 {
		t.Fatalf("after not equal to 5")
	}

}

func TestIncr(t *testing.T) {
	db := NewDB()
	db.Set("ho", []byte("12"))

	length, err := db.Incr("ho")
	val := db.Get("ho")
	valInt, err := strconv.Atoi(string(val))
	if err != nil {
		t.Fatalf(err.Error())
	}
	if length != valInt {
		t.Fatalf("before after not equal")
	}

	if valInt != 13 {
		t.Fatalf("after not equal to 13")
	}

	length, err = db.Incr("ho")
	val = db.Get("ho")
	valInt, err = strconv.Atoi(string(val))
	if err != nil {
		t.Fatalf(err.Error())
	}
	if length != valInt {
		t.Fatalf("before after not equal")
	}

	if valInt != 14 {
		t.Fatalf("after not equal to 14")
	}

}

func TestIncrby(t *testing.T) {
	db := NewDB()
	db.Set("ho", []byte("12"))

	length, err := db.Incrby("ho", 4)
	val := db.Get("ho")
	valInt, err := strconv.Atoi(string(val))
	if err != nil {
		t.Fatalf(err.Error())
	}
	if length != valInt {
		t.Fatalf("before after not equal")
	}

	if valInt != 16 {
		t.Fatalf("after not equal to 16")
	}

	length, err = db.Incrby("ho", 10)
	val = db.Get("ho")
	valInt, err = strconv.Atoi(string(val))
	if err != nil {
		t.Fatalf(err.Error())
	}
	if length != valInt {
		t.Fatalf("before after not equal")
	}

	if valInt != 26 {
		t.Fatalf("after not equal to 26")
	}

}

func TestGetrange(t *testing.T) {
	db := NewDB()

	db.Set("ho", []byte("This is a sentence."))

	// Good range
	val := string(db.Getrange("ho", 6, 12))
	if val != "s a se" {
		t.Fatal(val)
	}

	// From the end
	val = string(db.Getrange("ho", 0, -1))
	if val != "This is a sentence." {
		t.Fatal(val)
	}

	val = string(db.Getrange("ho", 0, -3))
	if val != "This is a sentenc" {
		t.Fatal(val)
	}

	val = string(db.Getrange("ho", -5, -3))
	if val != "nc" {
		t.Fatal(val)
	}

	// Out of range
	val = string(db.Getrange("ho", 3, 1))
	if val != "" {
		t.Fatal(val)
	}

	val = string(db.Getrange("ho", -3, -5))
	if val != "" {
		t.Fatal(val)
	}

	val = string(db.Getrange("ho", 10, 100))
	if val != "sentence." {
		t.Fatal(val)
	}
}
