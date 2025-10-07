package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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
			log.Println("accept error:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Accepted connection from", conn.RemoteAddr())

	resp := NewResp(conn)
	writer := NewWriter(conn)
	for {
		v, err := resp.Read()
		if err == io.EOF {
			fmt.Println("Connection closed by", conn.RemoteAddr())
			return
		}
		if err != nil {
			log.Println("Error reading from connection:", err)
			return
		}

		cmd, args, err := parseCommand(v)
		if err != nil {
			log.Println("Error parsing the command:", err)
			return
		}

		fmt.Printf("Received: %s %s\n", cmd, strings.Join(args, " "))
		writer.Write(Value{typ: "string", str: "OK"})
	}
}
