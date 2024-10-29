package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
}

// Returns a pointer type.
// Initializing and returning a new Server object
func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

// This allows us to add methods specific to the Server type variables.
// Enables the WS connection, to true in the Server object.
func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("New incoming connection from the client", ws.RemoteAddr())
	s.conns[ws] = true
	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {

	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)

		if err != nil {

			if err == io.EOF {
				break
			}

			fmt.Println("Read error: ", err)
			continue
		}

		msg := buf[:n]
		fmt.Println(string(msg))

		s.broadCast(msg)

		// ws.Write([]byte("Thank you for the message!!"))

	}
}

func main() {
	server := NewServer()
	http.Handle("/ws", websocket.Handler(server.handleWS))
	http.Handle("/orderbook-feed", websocket.Handler(server.handleWSOrderbook))
	http.ListenAndServe("localhost:3000", nil)
}

func (s *Server) broadCast(b []byte) {
	for ws := range s.conns {
		// A Go Routine to broadcast the message to the all other clients
		go func(wsocket *websocket.Conn) {
			_, err := ws.Write(b)
			if err != nil {
				fmt.Println("Write error = ", err)
			}
		}(ws)
	}
}

func (s *Server) handleWSOrderbook(ws *websocket.Conn) {
	fmt.Println("New incoming connection from the client to the orderbook feed", ws.RemoteAddr())

	for {
		payload := fmt.Sprintf("Orderbook data = %d\n", time.Now().UnixNano())
		ws.Write([]byte(payload))
		time.Sleep(time.Second * 2)
	}

}
