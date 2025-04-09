package server

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	_ "github.com/glebarez/go-sqlite"
)

type Server struct {
	db       *sql.DB
	wg       sync.WaitGroup
	quit     chan interface{}
	listener net.Listener
}

func NewServer(addr string) (*Server, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on address: %v", addr)
	}

	db, err := sql.Open("sqlite", "./db/lotus.db")
	if err != nil {
		log.Fatalf("couldn't open the database: %q", err)
	}

	return &Server{
		db:       db,
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

	user := s.receiveUser(conn)

	tickerExp := time.NewTicker(10 * time.Second)
	defer tickerExp.Stop()
	tickerLevel := time.NewTicker(30 * time.Second)
	defer tickerLevel.Stop()

	for {
		select {
		case <-tickerExp.C:
			user.exp += 10
			_, err := conn.Write([]byte(fmt.Sprintf("%s, you've gained 10 EXP. you have now %d EXP.\n", user.name, user.exp)))
			if err != nil {
				s.handleDisconnect(user)
				return
			}
		case <-tickerLevel.C:
			_, err := conn.Write([]byte(fmt.Sprintf("%s, you are level %d.\n", user.name, user.calculateLevel())))
			if err != nil {
				s.handleDisconnect(user)
				return
			}
		}
	}
}

func (s *Server) Stop() {
	close(s.quit)
	s.listener.Close()
	s.db.Close()
	s.wg.Wait()
}

func (s *Server) handleDisconnect(user *User) {
	log.Println("user disconnected:", user.name)
	if err := s.updateUser(user); err != nil {
		log.Println("failed updating user:", err)
	}
}

func (s *Server) receiveUser(conn net.Conn) *User {
	buf := make([]byte, 1024)
	size, err := conn.Read(buf)
	if err != nil {
		log.Println(err)
	}
	data := buf[:size]

	name := strings.Trim(string(data), "\n ")
	user := s.getUserOrCreate(name)

	log.Printf("User logged in: %v\n", user.name)

	conn.Write([]byte(fmt.Sprintf("Welcome, %s. You are IN the game.", user.name)))

	return user
}

func (s *Server) getUserOrCreate(name string) *User {
	row := s.db.QueryRow("select * from user where name = ?", name)

	var user User

	if err := row.Scan(&user.id, &user.name, &user.exp); err != nil {
		if err == sql.ErrNoRows {
			user = User{
				name: name,
				exp:  0,
			}
			s.addUser(&user)
		} else {
			log.Fatal(err)
		}
	}

	return &user
}

func (s *Server) addUser(user *User) (int64, error) {
	result, err := s.db.Exec("INSERT INTO user (name, exp) VALUES (?, ?)", user.name, user.exp)
	if err != nil {
		return 0, fmt.Errorf("addUser: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addUser: %v", err)
	}
	return id, nil
}

func (s *Server) updateUser(user *User) error {
	if _, err := s.db.Exec("UPDATE user SET exp = ? where name = ?;", user.exp, user.name); err != nil {
		return fmt.Errorf("updateUser: %v", err)
	}
	return nil
}
