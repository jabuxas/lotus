package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jabuxas/lotus/internal/server"
)

func main() {
	sv, err := server.NewServer(":6969")
	if err != nil {
		log.Println(err)
	}

	sv.Serve()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("shutting down..")
	sv.Stop()
}
