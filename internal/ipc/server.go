package ipc

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type HandlerFunc func(Request) Response

type Server struct {
	socketPath string
	handler    HandlerFunc

	listener net.Listener
}

func NewServer(
	socketPath string,
	handler HandlerFunc,
) *Server {
	return &Server{
		socketPath: socketPath,
		handler:    handler,
	}
}

func (s *Server) Start() error {
	if err := os.RemoveAll(s.socketPath); err != nil {
		return err
	}

	listener, err := net.Listen("unix", s.socketPath)
	if err != nil {
		return fmt.Errorf("listen unix socket: %w", err)
	}

	if err := os.Chmod(s.socketPath, 0660); err != nil {
		listener.Close()

		return err
	}

	s.listener = listener

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) Close() error {
	if s.listener == nil {
		return nil
	}

	return s.listener.Close()
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	var req Request

	if err := json.NewDecoder(conn).Decode(&req); err != nil {
		_ = json.NewEncoder(conn).Encode(Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	resp := s.handler(req)

	_ = json.NewEncoder(conn).Encode(resp)
}
