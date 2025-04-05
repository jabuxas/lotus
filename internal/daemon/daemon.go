package daemon

import (
	"log"
	"net"
)

const (
	Addr  = ":6969"
	Host  = "localhost"
	Proto = "tcp"
)

func StartDaemon() {
	ln, err := createListener()
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)

		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		buf := make([]byte, 1024)
		size, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}
		data := buf[:size]
		log.Printf("Received: %v\n", string(data))
		conn.Write(data)
	}
}

func createListener() (net.Listener, error) {
	ln, err := net.Listen("tcp", Addr)
	return ln, err
}
