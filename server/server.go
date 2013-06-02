package main

import (
	// "bytes"
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	// "os"
	"github.com/nicolaspaton/rhube"
)

func main() {
	fmt.Println("Rhube Server started.")
	ln, err := net.Listen("tcp", ":6380")
	if err != nil {
		panic(err)
	}
	db := rhube.NewDB()
	db.Set("toto", []byte("hey"))
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
			continue
		}
		go handleConnection(conn, db)
	}
}

func handleConnection(c net.Conn, db *rhube.DB) {

	bufIn := bufio.NewReader(c)
	bufOut := bufio.NewWriter(c)
	buf := bufio.NewReadWriter(bufIn, bufOut)

	var line []byte
	var err error
	partsExp, partsGot := 0, 0
	bytesExp, bytesGot := 0, 0
	binIn := false
	// bstringBuf := make([]byte, 1024)
	bString := make([][]byte, 0)
	for {
		if !binIn {
			line, _, err = buf.ReadLine()
			if err != nil {
				if err == io.EOF {
					fmt.Println("Connection closed.")
					partsExp, partsGot = 0, 0
					bytesExp, bytesGot = 0, 0
					line = []byte("")
					return
				} else {
					panic(err)
				}
			}

			unit := string(line)
			if unit == "" {
				continue
			}

			switch unit[:1] {
			case "*":
				partsExp, err = strconv.Atoi(unit[1:])
				if err != nil {
					panic(err)
				}
				// fmt.Println("multiBulkCount:", partsExp)
			case "$", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
				partsGot++
				if unit[:1] == "$" {
					bytesExp, err = strconv.Atoi(unit[1:])
				} else {
					bytesExp, err = strconv.Atoi(unit)
				}

				if err != nil {
					panic(err)
				}
				// fmt.Println("bytesExp:", bytesExp)
				binIn = true
			}
		} else {
			var tmpstr []byte
			var tmpGot int
			for bytesExp > bytesGot {
				line, err = buf.ReadSlice('\n')
				if err != nil {
					panic(err)
				}
				tmpGot = len(line)
				bytesGot += tmpGot + 1
				line = append(line, '\n')
				tmpstr = append(tmpstr, line...)
			}
			bytesGot = 0
			if tmpGot == 1 {
				bString = append(bString, tmpstr[:tmpGot-1])
			} else if len(tmpstr) > 0 {
				bString = append(bString, tmpstr[:tmpGot-2])
			}

			binIn = false
			if partsExp == len(bString) {
				partsGot, bytesGot = 0, 0
				result := ProcessRequest(bString, db)
				bString = make([][]byte, 0)
				buf.WriteString(result)
				buf.Flush()
			}
		}
	}

	return
	io.Copy(c, c)
}

func ProcessRequest(parts [][]byte, db *rhube.DB) string {
	// for _, part := range parts {
	// 	fmt.Print(string(part), " ")
	// }
	// fmt.Print("\n")

	switch strings.ToLower(string(parts[0])) {
	case "info":
		return ProcessInfo(db)
	case "get":
		val := db.Get(string(parts[1]))
		// fmt.Println(">", string(val))
		return "$" + strconv.Itoa(len(val)) + "\r\n" + string(val) + "\r\n"
	case "set":
		db.Set(string(parts[1]), parts[2])
		// fmt.Println("> +OK")
		return "+OK\r\n"
	}
	// fmt.Printf(strings.ToLower(string(parts[0])))
	// fmt.Printf("Here:%s", string(parts[0]), parts, len(parts))
	return "-Err cmd unprocessable"
}

func ProcessInfo(db *rhube.DB) string {
	return "$11\r\nhello:world\r\n"
}
