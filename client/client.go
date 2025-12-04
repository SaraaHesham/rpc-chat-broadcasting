package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strings"
)

type Message struct {
	User    string
	Content string
}

type Empty struct{}

type JoinArgs struct {
	UserID     string
	ClientAddr string
}

type Client struct{}

func (c *Client) Receive(msg Message, reply *Empty) error {
	fmt.Printf("\n[%s] %s\n> ", msg.User, msg.Content)
	return nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	clientRPC := new(Client)
	rpc.Register(clientRPC)

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal(err)
	}
	clientAddr := listener.Addr().String()

	go func() {
		for {
			conn, _ := listener.Accept()
			go rpc.ServeConn(conn)
		}
	}()

	server, err := rpc.Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal("Cannot connect to server:", err)
	}

	var reply Empty
	err = server.Call("ChatServer.Join",
		JoinArgs{UserID: name, ClientAddr: clientAddr},
		&reply)
	if err != nil {
		log.Fatal("Join failed:", err)
	}

	fmt.Println("Connected! Start typing...")

	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if text == "exit" {
			fmt.Println("Bye!")
			break
		}

		server.Call("ChatServer.SendMessage",
			Message{User: name, Content: text},
			&reply)
	}
}
