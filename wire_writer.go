package rhube

import (
	"bufio"
	"io"
	"strconv"
)

type WireWriter struct {
	w   io.Writer
	ww  *bufio.Writer
	buf []byte
}

const WireSep = "\r\n"

func NewWireWriter(r io.Writer) *WireWriter {
	return &WireWriter{r, bufio.NewWriter(r), make([]byte, bufSize)}
}

func (w *WireWriter) WriteStatusReply(s string) (err error) {
	_, err = w.ww.WriteString("+" + s + WireSep)
	if err != nil {
		return err
	}
	w.ww.Flush()
	return err
}

func (w *WireWriter) WriteErrorReply(e string) (err error) {
	_, err = w.ww.WriteString("-" + e + WireSep)
	if err != nil {
		return err
	}
	w.ww.Flush()
	return err
}

func (w *WireWriter) WriteIntReply(i int) (err error) {
	_, err = w.ww.WriteString(":" + strconv.Itoa(i) + WireSep)
	if err != nil {
		return err
	}
	w.ww.Flush()
	return err
}

func (w *WireWriter) WriteBoolReply(truth bool) (err error) {
	if truth {
		err = w.WriteIntReply(1)
	} else {
		err = w.WriteIntReply(0)
	}
	return err
}

func (w *WireWriter) WriteBulkReply(data []byte) (err error) {
	err = w.writeKey(string(data))
	if err != nil {
		return err
	}
	w.ww.Flush()
	return err
}

func (w *WireWriter) WriteMultiBulkReply(data [][]byte) (err error) {
	_, err = w.ww.WriteString("*" + strconv.Itoa(len(data)) + WireSep)
	if err != nil {
		return err
	}

	for _, chunk := range data {
		err = w.writeKey(string(chunk))
		if err != nil {
			return err
		}
	}

	w.ww.Flush()
	return err
}

// TODO: dry up, WriteMultiBulkReply and WriteStringMultiBulkReply are awefully alike
func (w *WireWriter) WriteStringMultiBulkReply(data []string) (err error) {
	_, err = w.ww.WriteString("*" + strconv.Itoa(len(data)) + WireSep)
	if err != nil {
		return err
	}

	for _, chunk := range data {
		err = w.writeKey(string(chunk))
		if err != nil {
			return err
		}
	}

	w.ww.Flush()
	return err
}

// TODO: dry up, WriteMultiBulkReply and WriteStringMultiBulkReply are awefully alike
func (w *WireWriter) WriteHashMultiBulkReply(data map[string]string) (err error) {
	_, err = w.ww.WriteString("*" + strconv.Itoa(len(data)) + WireSep)
	if err != nil {
		return err
	}

	for key, value := range data {
		err = w.writeKey(key)
		if err != nil {
			return err
		}
		err = w.writeKey(value)
		if err != nil {
			return err
		}
	}

	w.ww.Flush()
	return err
}

func (w *WireWriter) writeKey(key string) error {
	if key == "" {
		_, err := w.ww.WriteString("$-1\r\n")
		if err != nil {
			return err
		}
	}

	_, err := w.ww.WriteString("$" + strconv.Itoa(len(key)) + WireSep)
	if err != nil {
		return err
	}

	_, err = w.ww.Write([]byte(key))
	if err != nil {
		return err
	}

	_, err = w.ww.Write([]byte(WireSep))
	if err != nil {
		return err
	}
	return err
}
