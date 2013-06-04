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

func (w *WireWriter) WriteBoolReply(b bool) (err error) {
	if b {
		err = w.WriteIntReply(1)
	} else {
		err = w.WriteIntReply(0)
	}
	return err
}

func (w *WireWriter) WriteBulkReply(data []byte) (err error) {
	l := len(data)
	_, err = w.ww.WriteString("$" + strconv.Itoa(l) + WireSep)
	if err != nil {
		return err
	}
	_, err = w.ww.Write(data)
	if err != nil {
		return err
	}
	_, err = w.ww.Write([]byte(WireSep))
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
		_, err = w.ww.WriteString("$" + strconv.Itoa(len(chunk)) + WireSep)
		if err != nil {
			return err
		}

		_, err = w.ww.Write(chunk)
		if err != nil {
			return err
		}
		_, err = w.ww.Write([]byte(WireSep))
		if err != nil {
			return err
		}
	}
	w.ww.Flush()
	return err
}
