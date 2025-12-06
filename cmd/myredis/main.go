package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/AndrewSukhobok95/go-build-my-own-redis/internal/commands"
	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/persistence"
	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/server"
	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/storage"
)

func main() {
	fmt.Println("Starting MyRedis server...")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	storage := storage.NewKV()
	aof, err := persistence.NewAOF("appendonly.aof")
	if err != nil {
		panic(fmt.Sprintf("Can't open AOF file: %v", err))
	}
	server := server.New("6380", storage, aof, time.Duration(5)*time.Second)

	go server.Start()

	<-stop
	fmt.Println("\nShutting down gracefully...")
	server.Shutdown()
}
