package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("Starting MyRedis server...")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	storage := NewStorage()
	server := NewServer("6380", storage, time.Duration(5)*time.Second)

	go server.Start()

	<-stop
	fmt.Println("\nShutting down gracefully...")
	server.Shutdown()
}
