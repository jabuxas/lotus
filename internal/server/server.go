package server

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type Server struct {
	wg       sync.WaitGroup
	quit     chan interface{}
	listener net.Listener
}

func NewServer(addr string) (*Server, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on address: %v", addr)
	}

	return &Server{
		listener: ln,
		quit:     make(chan interface{}),
	}, nil
}

func (s *Server) Serve() {
	defer s.wg.Done()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.quit:
				return
			default:
				log.Println(err)
			}
		} else {
			s.wg.Add(1)
			go func() {
				s.handleConnection(conn)
				s.wg.Done()
			}()
		}
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	user := receiveUser(conn)

	tickerExp := time.NewTicker(10 * time.Second)
	defer tickerExp.Stop()
	tickerLevel := time.NewTicker(30 * time.Second)
	defer tickerLevel.Stop()

	for {
		select {
		case <-tickerExp.C:
			user.exp += 10
			_, err := conn.Write([]byte(fmt.Sprintf("%s, you've gained 10 EXP. you have now %d EXP.", user.name, user.exp)))
			if err != nil {
				log.Println("client disconnected: ", err)
				return
			}
		case <-tickerLevel.C:
			_, err := conn.Write([]byte(fmt.Sprintf("%s, you are level %d", user.name, user.calculateLevel())))
			if err != nil {
				log.Println("client disconnected: ", err)
				return
			}
		}
	}

}

func (s *Server) Stop() {
	close(s.quit)
	s.listener.Close()
	s.wg.Wait()
}

func receiveUser(conn net.Conn) *User {
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

	return user
}
