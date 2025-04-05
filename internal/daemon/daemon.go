package daemon

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

const (
	Addr  = ":6969"
	Host  = "localhost"
	Proto = "tcp"
)

func StartDaemon() {
	ln, err := net.Listen("tcp", Addr)
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

	user := receiveUser(conn)

	for {
		time.Sleep(time.Second * 5)
		conn.Write([]byte(fmt.Sprintf("%s, you are still in the game.", user.name)))
	}
}

type User struct {
	name string
}

func receiveUser(conn net.Conn) User {
	buf := make([]byte, 1024)
	size, err := conn.Read(buf)
	if err != nil {
		log.Println(err)
	}
	data := buf[:size]

	user := &User{
		name: strings.Trim(string(data), "\n "),
	}

	log.Printf("New user: %v\n", user.name)

	conn.Write([]byte(fmt.Sprintf("Welcome, %s. You are IN the game.", user.name)))

	return *user
}
