package main

import (
	"fmt"
	"github.com/nicolaspaton/rhube"
	"io"
	"net"
	"os"
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
		go handleConn(conn, db)
	}
}

func handleConn(c net.Conn, db *rhube.DB) {
	defer c.Close()
	r := rhube.NewWireReader(c)
	w := rhube.NewWireWriter(c)
	for {
		args, err := r.ReadCommand()
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed.")
				return
				return
			} else {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		// fmt.Println("args:", args)
		switch args[0] {
		case "info":
			w.WriteBulkReply([]byte("hello:world"))
		case "get":
			w.WriteBulkReply(db.Get(args[1]))
		case "set":
			w.WriteBoolReply(db.Set(args[1], []byte(args[2])))
		}
	}
}
