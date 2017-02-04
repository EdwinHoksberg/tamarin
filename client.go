package main

import (
	"net"
	"time"
)

type Client struct {
	connection      net.Conn
	connectionStart time.Time
	isNewConnection bool
}

func NewClient(client net.Conn) *Client {
	return &Client{client, time.Now(), true}
}

func (c *Client) getConnectionTime() time.Duration {
	return time.Since(c.connectionStart)
}
