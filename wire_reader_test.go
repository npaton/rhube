package rhube

import (
	// "log"
	// "strconv"
	"bytes"
	"testing"
)

func TestWireReader(t *testing.T) {
	buf := bytes.NewBuffer([]byte("*1\r\n$4\r\ninfo\r\n"))
	w := NewWireReader(buf)
	args, err := w.ReadCommand()
	if err != nil {
		t.Fatal(err.Error())
	}

	if args[0] != "info" {
		t.Fatalf("args[0] is not info:", args[0])
	}

	buf = bytes.NewBuffer([]byte("*2\r\n$3\r\nget\r\n$5\r\nfooba\r\n"))
	w = NewWireReader(buf)
	args, err = w.ReadCommand()
	if err != nil {
		t.Fatal(err.Error())
	}

	if args[0] != "get" || args[1] != "fooba" {
		t.Fatalf("args is not get fooba:", args)
	}

	buf = bytes.NewBuffer([]byte("*1\r\n$4\r\ninfo\r\n*2\r\n$3\r\nget\r\n$5\r\nfooba\r\n*1\r\n$4\r\ninfo\r\n"))
	w = NewWireReader(buf)
	args, err = w.ReadCommand()
	if err != nil {
		t.Fatal(err.Error())
	}

	if args[0] != "info" {
		t.Fatalf("args is not get fooba:", args)
	}

	args, err = w.ReadCommand()
	if err != nil {
		t.Fatal(err.Error())
	}

	if args[0] != "get" || args[1] != "fooba" {
		t.Fatalf("args is not get fooba:", args)
	}

	args, err = w.ReadCommand()
	if err != nil {
		t.Fatal(err.Error())
	}

	if args[0] != "info" {
		t.Fatalf("args is not get fooba:", args)
	}

	// EOF
	args, err = w.ReadCommand()
	if err == nil {
		t.Fatal("err should not be nil!")
	}

}


