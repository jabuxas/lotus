package client

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/jabuxas/lotus/internal/server"
)

func Client() {
	conn, err := net.Dial(server.Proto, fmt.Sprintf("%s%s", server.Host, server.Addr))
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	go sendHello(conn)

	for {
		buf := make([]byte, 1024)

		size, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Fatalf("client can't connect to the daemon: %q", err)
			}
			log.Println(err)
		}

		msg := string(buf[:size])
		fmt.Println(msg)
	}

}

func sendHello(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Type your username: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		log.Println(err)
	}

	buf := make([]byte, 1024)
	size, err := conn.Write([]byte(name))

	if err != nil {
		log.Println(err)
		return
	}

	data := buf[:size]
	conn.Read(data)

	fmt.Printf("%v\n", string(data))
}
