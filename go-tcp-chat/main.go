package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

type Message struct {
	ClientID string
	Message  string
}

type Client struct {
	Username    string
	ClientId    string
	Conn        net.Conn
	IsConnected bool
}

var (
	Clients   = make(map[string]Client)
	clientsMu sync.RWMutex
)

func GenerateClientIDHex(byteSize int) (string, error) {
	b := make([]byte, byteSize)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil

}

func broadCastMessages(ctx context.Context, broadcastReq chan Message, wg *sync.WaitGroup) {

	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("received signal to close the connection")
			return
		case message := <-broadcastReq:
			fmt.Println("sending a message from", message.ClientID, "to all the connected clients")
			clientsMu.RLock()
			for clientID, connection := range Clients {
				_, err := connection.Conn.Write([]byte(message.Message))
				if err != nil {
					fmt.Printf("Server write error: %v to client %s", err, clientID)
					continue
				}

			}
			clientsMu.RUnlock()
		}
	}
	return
}

func handleConnection(ctx context.Context, clientID string, broadCastChannel chan Message, wg *sync.WaitGroup) {

	defer wg.Done()
	fmt.Println("waiting to acquire lock to the client list")
	clientsMu.Lock()
	fmt.Println("acquired lock to the client list")

	if _, found := Clients[clientID]; !found {
		fmt.Println("failed to find an active connection for the client id :", clientID)
		return
	}

	client := Clients[clientID]

	defer func() {
		client.Conn.Close() // Free the network port
		clientsMu.Lock()
		delete(Clients, clientID) // Remove from global map
		clientsMu.Unlock()

	}()

	_, err := client.Conn.Write([]byte("Hello User, Welcome to The Chat Session, What should I call you ?    \n"))

	reader := bufio.NewReader(client.Conn)
	username, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Read error: %v", err)
		fmt.Println("closing the connection")
		return
	}
	if username != "" {
		usernameParts := strings.Split(username, " ")
		if len(usernameParts) > 0 {
			client.Username = strings.TrimRight(usernameParts[0], "\r\n")
			Clients[clientID] = client
		}

	}
	if username == "" {
		client.Username = client.ClientId[:5]
		Clients[clientID] = client

	}
	clientsMu.Unlock()
	fmt.Println("released lock from the client list")

	_, err = client.Conn.Write([]byte(fmt.Sprintf("Start typing %sto send message to everyone\n", client.Username)))

	if err != nil {
		fmt.Printf("Read error: %v", err)
		fmt.Println("closing the connection")
		return
	}
	for {

		// Handle connection specific I/O here...
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Read error: %v", err)
			return
		}

		broadCastChannel <- Message{
			ClientID: clientID,
			Message:  fmt.Sprintf("%s:%s", client.Username, message),
		}
	}
}

func startServer(ctx context.Context, portNumber string) {

	// listner
	lc := net.ListenConfig{}
	listener, err := lc.Listen(ctx, "tcp", portNumber)
	if err != nil {
		fmt.Println("failed to start a server on the port ", portNumber, "err:", err)
		return
	}

	fmt.Println("server is listening on :", portNumber)

	go func() {
		<-ctx.Done()
		fmt.Println("Context canceled, closing listener...")
		// Closing the listener forces Accept() to unblock immediately
		listener.Close()

		clientsMu.Lock()
		for _, client := range Clients {
			// Closing the connection forces any blocked bufio.Readers to instantly return an error
			client.Conn.Close()
		}
		clientsMu.Unlock()
	}()

	// client connection handler

	// this is a broadcast channel, that sends a message from
	// one client to all connected clients on this server.

	broadCastChannel := make(chan Message, 1000)
	var wg sync.WaitGroup

	wg.Add(1)
	go broadCastMessages(ctx, broadCastChannel, &wg)
	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				fmt.Println("Accept loop terminated gracefully.")
				break
			}
			fmt.Println("failed to accept an incoming connection request", err)
			continue
		}

		hexID, err := GenerateClientIDHex(16)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			fmt.Println("failed to handle a  connection request")
			continue
		}
		fmt.Println("waiting to acquire lock to the client list")
		clientsMu.Lock()
		fmt.Println("acquired lock to the client list")
		// add all new connections to a global pool
		Clients[hexID] = Client{
			ClientId:    hexID,
			Conn:        conn,
			IsConnected: true,
		}
		clientsMu.Unlock()
		fmt.Println("released lock for the client list")
		wg.Add(1)
		go handleConnection(ctx, hexID, broadCastChannel, &wg)

	}

	fmt.Println("waiting for all the clients to close")
	ctx.Done()
	wg.Wait()
	fmt.Println("all clients terminated")

}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	startServer(ctx, ":6000")
}
