package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/engine"
	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/persistence"
	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/resp"
	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/storage"
)

func handleConnection(conn net.Conn, storage *storage.KV, aof *persistence.AOF, shutdown <-chan struct{}) {
	defer conn.Close()
	fmt.Println("Accepted connection from", conn.RemoteAddr())

	ctx := engine.NewCommandContext(storage, aof)

	respReader := resp.NewReader(conn)
	respWriter := resp.NewWriter(conn)
	for {
		conn.SetReadDeadline(time.Now().Add(time.Second))
		v, err := respReader.Read()
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
			respWriter.Write(resp.NewErrorValue("ERR invalid command"))
			continue
		}

		cmd, args, err := resp.ParseCommand(v)
		if err != nil {
			log.Println("Error parsing the command:", err)
			respWriter.Write(resp.NewErrorValue("ERR invalid command"))
			continue
		}

		fmt.Printf("Received: %s %s\n", cmd, strings.Join(args, " "))
		answerValue := engine.DispatchCommand(ctx, cmd, args)
		respWriter.Write(answerValue)
	}
}

type Server struct {
	listener        net.Listener
	storage         *storage.KV
	aof             *persistence.AOF
	wg              sync.WaitGroup
	shutdown        chan struct{}
	addr            string
	cleanupInterval time.Duration
}

func New(addr string, storage *storage.KV, aof *persistence.AOF, cleanupInterval time.Duration) *Server {
	return &Server{storage: storage, aof: aof, addr: addr, shutdown: make(chan struct{}, 1), cleanupInterval: cleanupInterval}
}

func (s *Server) Start() {
	var err error
	s.listener, err = net.Listen("tcp", fmt.Sprintf(":%s", s.addr))
	if err != nil {
		log.Fatalln(err)
	}

	s.ReplayAOF()
	go s.FlushAOF()
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
			handleConnection(conn, s.storage, s.aof, s.shutdown)
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

func (s *Server) ReplayAOF() {
	log.Println("Starting Replay AOF")
	ctx := engine.NewCommandContext(s.storage, s.aof)
	ctx.StartReplay()

	cmdCh := make(chan persistence.ReplayCommand, 10)
	go func() {
		err := s.aof.Load(cmdCh)
		if err != nil {
			log.Printf("Error while AOF Replaying: %v", err)
		}
	}()
	for cmd := range cmdCh {
		log.Printf("Dispatching %v %v", cmd.Name, cmd.Args)
		engine.DispatchCommand(ctx, cmd.Name, cmd.Args)
	}
	log.Println("Replay AOF finished")
}

func (s *Server) FlushAOF() {
	for {
		select {
		case <-s.shutdown:
			return
		case <-time.After(time.Second):
			if err := s.aof.Flush(); err != nil {
				log.Println("Flush error:", err)
			}
		}
	}
}
