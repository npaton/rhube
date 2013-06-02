package main

import (
	"fmt"
	"io"
	"net"
)

func main() {
	fmt.Println("Rhube Server started.")
	ln, err := net.Listen("tcp", ":6380")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
	io.Copy(c, c)
}
