package client

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func StartClient(host, port string) {
	srv := fmt.Sprint(host, port)
	conn, err := net.Dial("tcp", srv)
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
				log.Println("server has closed the connection.")
				return
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
