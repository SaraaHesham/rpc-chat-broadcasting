package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
)

type Message struct {
	User    string
	Content string
}

type JoinArgs struct {
	UserID     string
	ClientAddr string
}

type Empty struct{}

type ChatServer struct {
	mu      sync.Mutex
	clients map[string]string

	broadcastChan chan Message
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		clients:       make(map[string]string),
		broadcastChan: make(chan Message, 100),
	}
}

func (s *ChatServer) Join(args JoinArgs, reply *Empty) error {
	s.mu.Lock()
	s.clients[args.UserID] = args.ClientAddr
	s.mu.Unlock()

	fmt.Println("Client Joined:", args.UserID)

	// send join message to broadcaster
	s.broadcastChan <- Message{
		User:    "SYSTEM",
		Content: fmt.Sprintf("User %s joined", args.UserID),
	}

	return nil
}

func (s *ChatServer) SendMessage(msg Message, reply *Empty) error {
	fmt.Println("Received from", msg.User, ":", msg.Content)

	// push message into broadcast channel
	s.broadcastChan <- msg
	return nil
}

func (s *ChatServer) broadcaster() {
	for msg := range s.broadcastChan {

		s.mu.Lock()
		clientsCopy := make(map[string]string)
		for id, addr := range s.clients {
			clientsCopy[id] = addr
		}
		s.mu.Unlock()

		for user, addr := range clientsCopy {
			if user == msg.User {
				continue
			}

			go func(clientAddr string) {
				client, err := rpc.Dial("tcp", clientAddr)
				if err != nil {
					fmt.Println("Broadcast error:", err)
					return
				}
				defer client.Close()

				var r Empty
				client.Call("Client.Receive", msg, &r)
			}(addr)
		}
	}
}

func main() {
	server := NewChatServer()
	rpc.Register(server)

	// start broadcast goroutine
	go server.broadcaster()

	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server running at 1234 ...")

	for {
		conn, _ := l.Accept()
		go rpc.ServeConn(conn)
	}
}
