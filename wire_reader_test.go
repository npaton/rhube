package rhube

import (
	// "log"
	"bytes"
	"strconv"
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

var bigVal string

const bigValSize = 10000000 // ~10MB

func init() {
	bigVal = string(make([]byte, bigValSize))
}

func TestWireReaderLarge(t *testing.T) {
	in := "*3\r\n$3\r\nset\r\n$3\r\nfoo\r\n$" + strconv.Itoa(bigValSize) + "\r\n" + bigVal + "\r\n"
	// t.Error(in)
	buf := bytes.NewBuffer([]byte(in))
	w := NewWireReader(buf)

	args, err := w.ReadCommand()
	if err != nil {
		t.Fatal(err.Error())
	}

	if args[0] != "set" || args[1] != "foo" || args[2] != bigVal {
		t.Fatalf("args is not get bigVal:", args, len(args[2]))
	}

	// EOF
	args, err = w.ReadCommand()
	if err == nil {
		t.Fatal("err should not be nil!")
	}

}
