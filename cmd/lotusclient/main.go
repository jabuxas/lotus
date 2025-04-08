package main

import (
	"github.com/jabuxas/lotus/internal/client"
)

func main() {
	client.StartClient("localhost", ":6969")
}
