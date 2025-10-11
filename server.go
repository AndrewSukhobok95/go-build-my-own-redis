package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

func handleConnection(conn net.Conn, storage *Storage, shutdown <-chan struct{}) {
	defer conn.Close()
	fmt.Println("Accepted connection from", conn.RemoteAddr())

	resp := NewResp(conn)
	writer := NewWriter(conn)
	for {
		conn.SetReadDeadline(time.Now().Add(time.Second))
		v, err := resp.Read()
		nerr, ok := err.(net.Error)
		switch {
		case ok && nerr.Timeout():
			select {
			case <-shutdown:
				return
			default:
			}
		case err == io.EOF:
			log.Println("Connection closed by", conn.RemoteAddr())
			return
		case err != nil:
			log.Println("Error reading from connection:", err)
			return
		}

		cmd, args, err := parseCommand(v)
		if err != nil {
			log.Println("Error parsing the command:", err)
			return
		}

		fmt.Printf("Received: %s %s\n", cmd, strings.Join(args, " "))
		answerValue := handleCommand(cmd, args, storage)
		writer.Write(answerValue)
	}
}

type Server struct {
	listener net.Listener
	storage  *Storage
	wg       sync.WaitGroup
	shutdown chan struct{}
	addr     string
}

func NewServer(addr string, storage *Storage) *Server {
	return &Server{storage: storage, addr: addr, shutdown: make(chan struct{}, 1)}
}

func (s *Server) Start() {
	var err error
	s.listener, err = net.Listen("tcp", fmt.Sprintf(":%s", s.addr))
	if err != nil {
		log.Fatalln(err)
	}

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.shutdown:
				return
			default:
				log.Println("Accept error:", err)
				continue
			}
		}

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			handleConnection(conn, s.storage, s.shutdown)
		}()
	}
}

func (s *Server) Shutdown() {
	close(s.shutdown)
	err := s.listener.Close()
	if err != nil {
		log.Println("Can't close the listener:", err)
	}
	s.wg.Wait()
}
