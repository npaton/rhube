package rhube

import (
	// "log"
	"bytes"
	"strconv"
	"testing"
)

func TestWireWriterStatus(t *testing.T) {
	buf := bytes.NewBuffer([]byte(""))
	w := NewWireWriter(buf)
	err := w.WriteStatusReply("OK")
	if err != nil {
		t.Fatal(err.Error())
	}
	if buf.String() != "+OK\r\n" {
		t.Fatal("Should be '+OK'", buf.String())
	}
}

func TestWireWriterError(t *testing.T) {
	buf := bytes.NewBuffer([]byte(""))
	w := NewWireWriter(buf)
	err := w.WriteErrorReply("Error: Something went wrong")
	if err != nil {
		t.Fatal(err.Error())
	}
	if buf.String() != "-Error: Something went wrong\r\n" {
		t.Fatal("Should be '-Error: Something went wrong'", buf.String())
	}
}

func TestWireWriterInt(t *testing.T) {
	buf := bytes.NewBuffer([]byte(""))
	w := NewWireWriter(buf)
	err := w.WriteIntReply(12)
	if err != nil {
		t.Fatal(err.Error())
	}
	if buf.String() != ":12\r\n" {
		t.Fatal("Should be ':12\\r\\n'", buf.String())
	}

	buf.Reset()

	err = w.WriteIntReply(9234567890)
	if err != nil {
		t.Fatal(err.Error())
	}
	if buf.String() != ":9234567890\r\n" {
		t.Fatal("Should be ':9234567890\\r\\n'", buf.String())
	}
}

func TestWireWriterBulk(t *testing.T) {
	d := "Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor."
	buf := bytes.NewBuffer([]byte(""))
	w := NewWireWriter(buf)
	err := w.WriteBulkReply([]byte(d))
	if err != nil {
		t.Fatal(err.Error())
	}
	exp := "$" + strconv.Itoa(len(d)) + "\r\n" + d + "\r\n"
	if buf.String() != exp {
		t.Fatalf("Should be", exp, buf.String())
	}
}

func TestWireWriterMultiBulk(t *testing.T) {
	d := "Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor."
	buf := bytes.NewBuffer([]byte(""))
	w := NewWireWriter(buf)
	err := w.WriteMultiBulkReply([][]byte{[]byte(d), []byte(d)})
	if err != nil {
		t.Fatal(err.Error())
	}
	exp := "$" + strconv.Itoa(len(d)) + "\r\n" + d + "\r\n"

	if buf.String() != "*2\r\n"+exp+exp {
		t.Fatalf("Should be", exp, buf.String())
	}
}
