package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/commands"
	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/resp"
	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/storage"
)

func handleConnection(conn net.Conn, storage *storage.Storage, shutdown <-chan struct{}) {
	defer conn.Close()
	fmt.Println("Accepted connection from", conn.RemoteAddr())

	respParser := resp.NewResp(conn)
	writer := resp.NewWriter(conn)
	for {
		conn.SetReadDeadline(time.Now().Add(time.Second))
		v, err := respParser.Read()
		nerr, ok := err.(net.Error)
		switch {
		case ok && nerr.Timeout():
			select {
			case <-shutdown:
				return
			default:
				continue
			}
		case err == io.EOF:
			log.Println("Connection closed by", conn.RemoteAddr())
			return
		case err != nil:
			log.Println("Error reading from connection:", err)
			writer.Write(resp.NewErrorValue("ERR invalid command"))
			continue
		}

		cmd, args, err := resp.ParseCommand(v)
		if err != nil {
			log.Println("Error parsing the command:", err)
			writer.Write(resp.NewErrorValue("ERR invalid command"))
			continue
		}

		fmt.Printf("Received: %s %s\n", cmd, strings.Join(args, " "))
		answerValue := commands.HandleCommand(cmd, args, storage)
		writer.Write(answerValue)
	}
}

type Server struct {
	listener        net.Listener
	storage         *storage.Storage
	wg              sync.WaitGroup
	shutdown        chan struct{}
	addr            string
	cleanupInterval time.Duration
}

func NewServer(addr string, storage *storage.Storage, cleanupInterval time.Duration) *Server {
	return &Server{storage: storage, addr: addr, shutdown: make(chan struct{}, 1), cleanupInterval: cleanupInterval}
}

func (s *Server) Start() {
	var err error
	s.listener, err = net.Listen("tcp", fmt.Sprintf(":%s", s.addr))
	if err != nil {
		log.Fatalln(err)
	}

	go s.storage.Cleanup(s.cleanupInterval, s.shutdown)

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
