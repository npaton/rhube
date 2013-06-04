package rhube

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
)

const bufSize = 1000000 // ~1MB

type WireReader struct {
	r        io.Reader
	rr       *bufio.Reader
	buf      []byte
	dBuf     []byte
	partsExp int
	partsGot int
	bytesExp int
}

func NewWireReader(r io.Reader) *WireReader {
	return &WireReader{r, bufio.NewReader(r), make([]byte, bufSize), make([]byte, 2), 0, 0, 0}
}


func (w *WireReader) ReadCommand() (args []string, err error) {
	num, err := w.ReadNumOfArgs()
	if err != nil {
		return nil, err
	}

	args = make([]string, num)
	for i := 0; i < num; i++ {
		arg, err := w.readArg()
		if err != nil {
			return args, err
		}
		args[i] = string(arg)
		// fmt.Println(i, string(arg))
	}
	return args, nil
}

func (w *WireReader) readArg() ([]byte, error) {

	line, err := w.rr.ReadSlice('\n')
	if err != nil {
		return nil, err
	}

	if string(line[:1]) != "$" {
		return nil, fmt.Errorf("WireReader. ParseError: '$' expected and not found at byte", string(w.buf[:1]))
	}

	numOfBytes, err := strconv.Atoi(string(line[1 : len(line)-2]))
	if err != nil {
		return nil, err
	}

	if numOfBytes > len(w.buf) {
		return nil, errors.New("TODO: Number of bytes of Value is bigger than buffer, should increase w.buf size here.")
	}

	n, err := w.rr.Read(w.buf[:numOfBytes])
	if n == 0 || err != nil {
		return nil, err
	}

	// fmt.Println("->" + string(w.buf[:numOfBytes]) + "\n")

	// Skip next "\r\n"
	n, err = w.rr.Read(w.dBuf[:2])
	if n != 2 || err != nil || w.dBuf[0] != '\r' || w.dBuf[1] != '\n' {
		return nil, err
	}

	return w.buf[:numOfBytes], nil
}

func (w *WireReader) ReadNumOfArgs() (int, error) {
	n, err := w.rr.Read(w.buf[:1])
	if n != 1 || err != nil {
		return 0, err
	}

	if string(w.buf[:1]) != "*" {
		return 0, errors.New("WireReader. ParseError: '*' expected and not found at byte")
	}

	line, err := w.rr.ReadSlice('\n')
	if n == 0 || err != nil {
		return 0, err
	}

	numOfArgs, err := strconv.Atoi(string(line[:len(line)-2]))
	if n == 0 || err != nil {
		return 0, err
	}

	return numOfArgs, nil
}
