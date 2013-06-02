package rhube

import (
	// "log"
	"strconv"
	"testing"
)

func TestGetSet(t *testing.T) {
	db := NewDB()
	db.Set("ho", []byte("12"))
	val := db.Get("ho")
	if string(val) != "12" {
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

	length, err := db.Decrby("ho", 2)
	val := db.Get("ho")
	valInt, err := strconv.Atoi(string(val))
	if err != nil {
		t.Fatalf(err.Error())
	}
	if length != valInt {
		t.Fatalf("before after not equal")
	}
	if valInt != -2 {
		t.Fatalf("after not equal to -2")
	}

	db.Set("ho", []byte("12"))

	length, err = db.Decrby("ho", 2)
	val = db.Get("ho")
	valInt, err = strconv.Atoi(string(val))
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

	length, err := db.Incrby("ho", 4)
	val := db.Get("ho")
	valInt, err := strconv.Atoi(string(val))
	if err != nil {
		t.Fatalf(err.Error())
	}
	if length != valInt {
		t.Fatalf("before after not equal")
	}
	if valInt != 4 {
		t.Fatalf("after not equal to 16")
	}

	db.Set("ho", []byte("12"))

	length, err = db.Incrby("ho", 4)
	val = db.Get("ho")
	valInt, err = strconv.Atoi(string(val))
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

func TestIncrbyfloat(t *testing.T) {
	db := NewDB()

	length, err := db.Incrbyfloat("ho", 4.4)
	val := db.Get("ho")
	if err != nil || length != string(val) || string(val) != "4.4" {
		t.Fatalf("Incrybyfloat: should be 4.4, was", length)
	}

	db.Set("ho", []byte("12"))

	length, err = db.Incrbyfloat("ho", 4.4)
	val = db.Get("ho")
	if err != nil || length != string(val) || string(val) != "16.4" {
		t.Fatalf("Incrybyfloat: should be 16.4, was", length)
	}

	length, err = db.Incrbyfloat("ho", 10.12345678)
	val = db.Get("ho")
	if err != nil || length != string(val) || string(val) != "26.523456779999997" {
		t.Fatalf("Incrybyfloat: should be 26.523456779999997, was", length)
	}

	db.Set("ho", []byte("0"))

	length, err = db.Incrbyfloat("ho", 10.1234567890123456789)
	val = db.Get("ho")
	if err != nil || length != string(val) || string(val) != "10.123456789012346" {
		t.Fatalf("Incrybyfloat: should be 10.123456789012346, was", length)
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

func TestStrlen(t *testing.T) {
	db := NewDB()
	db.Set("ho", []byte("1234567"))
	length := db.Strlen("ho")
	if length != 7 {
		t.Fatalf("Wierd!", length)
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

func TestSetrange(t *testing.T) {
	db := NewDB()

	db.Set("ho", []byte("This is a sentence."))

	db.Setrange("ho", 10, []byte("future perfect."))
	val := string(db.Get("ho"))
	if val != "This is a future perfect." {
		t.Fatal(val)
	}

	db.Setrange("ho", 30, []byte("No?"))
	val = string(db.Get("ho"))
	if val != "This is a future perfect.\u0000\u0000\u0000\u0000\u0000No?" {
		t.Fatal(val)
	}
}

func TestGetset(t *testing.T) {
	db := NewDB()
	db.Set("ho", []byte("12"))
	val, _ := db.Getset("ho", []byte("boo"))
	if string(val) != "12" {
		t.Fatalf("Wierd!")
	}
	val, _ = db.Getset("ho", []byte("soo"))
	if string(val) != "boo" {
		t.Fatalf("Wierd!")
	}
	val, _ = db.Getset("ho", []byte("roo"))
	if string(val) != "soo" {
		t.Fatalf("Wierd!", string(val))
	}
}

func TestMget(t *testing.T) {
	db := NewDB()
	db.Set("ho", []byte("12"))
	db.Set("ha", []byte("14"))
	db.Set("hi", []byte("15"))

	val := db.Mget("ho", "ha")
	if string(val[0]) != "12" || string(val[1]) != "14" {
		t.Fatalf("Wierd!", string(val[0]), string(val[1]))
	}
	val = db.Mget("hi", "ho", "ha")
	if string(val[0]) != "15" || string(val[1]) != "12" || string(val[2]) != "14" {
		t.Fatalf("Wierd!", string(val[0]), string(val[1]), string(val[2]))
	}
	val = db.Mget("ho", "nono", "ho")
	if string(val[0]) != "12" || string(val[1]) != "" || string(val[2]) != "12" {
		t.Fatalf("Wierd!", string(val[0]), string(val[1]), string(val[2]))
	}
}

func TestMset(t *testing.T) {
	db := NewDB()
	db.Set("ho", []byte("12"))
	db.Set("ha", []byte("14"))
	db.Set("hi", []byte("15"))

	val := db.Mget("ho", "ha", "hi")
	if string(val[0]) != "12" || string(val[1]) != "14" || string(val[2]) != "15" {
		t.Fatalf("Wierd!", string(val[0]), string(val[1]), string(val[2]))
	}

	db.Mset("hi", "boo", "ho", "soo", "ha", "roo")

	val = db.Mget("hi", "ho", "ha")
	if string(val[0]) != "boo" || string(val[1]) != "soo" || string(val[2]) != "roo" {
		t.Fatalf("Wierd!", string(val[0]), string(val[1]), string(val[2]))
	}
}

func TestMsetnx(t *testing.T) {
	db := NewDB()
	db.Set("ho", nil)
	db.Set("ha", nil)
	db.Set("hi", nil)

	retVal := db.Msetnx("hi", "boo", "ho", "soo", "ha", "roo")

	val := db.Mget("ho", "ha", "hi")
	if retVal != true || string(val[0]) != "soo" || string(val[1]) != "roo" || string(val[2]) != "boo" {
		t.Fatalf("Wierd!", retVal, string(val[0]), string(val[1]), string(val[2]))
	}

	// db.Set("ho", nil)
	db.Set("ha", nil)
	db.Set("hi", nil)
	retVal = db.Msetnx("hi", "boo", "ho", "soo", "ha", "roo")

	val = db.Mget("hi", "ho", "ha")
	if retVal != false || string(val[0]) != "" || string(val[1]) != "soo" || string(val[2]) != "" {
		t.Fatalf("Wierd!", retVal, string(val[0]), string(val[1]), string(val[2]))
	}
}
