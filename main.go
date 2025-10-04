package main

import (
	"fmt"
	"net"
	"log"
	"io"
)

func main() {
    fmt.Println("Starting MyRedis server...")

	ln, err := net.Listen("tcp", ":6380")
	if err != nil {
		log.Fatalln(err)
	}
	for {
		conn, err := ln.Accept()
		
		if err != nil {
			log.Fatalln(err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
    defer conn.Close()

    fmt.Println("Accepted connection from", conn.RemoteAddr())

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err == io.EOF {
			fmt.Println("Connection closed by", conn.RemoteAddr())
			return
		}
		if err != nil {
			log.Println("Error reading from connection:", err)
			return
		}
		fmt.Printf("Received: %s\n", string(buf[:n]))
		conn.Write([]byte("+OK\r\n"))
	}
}