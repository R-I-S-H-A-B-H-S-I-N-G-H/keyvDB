package server

import (
	"fmt"
	"net"
)

// Server represents the TCP server
type Server struct {
	Address string
}

// NewServer initializes a new TCP server
func NewServer(address string) *Server {
	return &Server{Address: address}
}

// Start runs the TCP server
func (s *Server) Start(handleConnection func(net.Conn)) error {
	listener, err := net.Listen("tcp", s.Address)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	defer listener.Close()

	fmt.Println("Server started on", s.Address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn) // Handle each connection in a separate goroutine
	}
}
