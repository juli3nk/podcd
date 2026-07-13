package ipc

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type Client struct {
	socketPath string
	timeout    time.Duration
}

func NewClient(socketPath string) *Client {
	return &Client{
		socketPath: socketPath,
		timeout:    5 * time.Second,
	}
}

func (c *Client) Send(req Request) (*Response, error) {
	conn, err := net.DialTimeout(
		"unix",
		c.socketPath,
		c.timeout,
	)
	if err != nil {
		return nil, fmt.Errorf("connect to daemon: %w", err)
	}
	defer conn.Close()

	if err := conn.SetDeadline(
		time.Now().Add(c.timeout),
	); err != nil {
		return nil, err
	}

	if err := json.NewEncoder(conn).Encode(req); err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}

	var resp Response

	if err := json.NewDecoder(conn).Decode(&resp); err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if !resp.Success {
		return &resp, fmt.Errorf("response: %s", resp.Error)
	}

	return &resp, nil
}
