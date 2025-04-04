package daemon

import (
	"log"
	"net"
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

func handleConnection(c net.Conn) {
	defer c.Close()

	for {
		buf := make([]byte, 1024)
		size, err := c.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}
		data := buf[:size]
		log.Printf("Received: %v\n", string(data))
		c.Write(data)
	}
}

func createListener() (net.Listener, error) {
	ln, err := net.Listen("tcp", ":6969")
	return ln, err
}
