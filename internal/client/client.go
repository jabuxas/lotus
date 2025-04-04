package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func Client() {
	conn, err := net.Dial("tcp", "localhost:6969")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		buf := make([]byte, 1024)
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		size, err := conn.Write([]byte(line))
		if err != nil {
			log.Fatal(err)
			return
		}
		data := buf[:size]
		conn.Read(data)
		fmt.Println(string(data))
	}
}
