package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/server"
	"github.com/AndrewSukhobok95/go-build-my-own-redis/internal/storage"
)

func main() {
	fmt.Println("Starting MyRedis server...")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	storage := storage.NewStorage()
	server := server.NewServer("6380", storage, time.Duration(5)*time.Second)

	go server.Start()

	<-stop
	fmt.Println("\nShutting down gracefully...")
	server.Shutdown()
}
