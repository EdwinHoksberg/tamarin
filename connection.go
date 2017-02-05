package main

import (
	"net"
	"strconv"
	"time"
)

type handleClientFunc func(*Client)

type Connection struct {
	server net.Listener
}

func (c *Connection) start(hostname string, portNumber int) bool {
	addr := hostname + ":" + strconv.Itoa(portNumber)

	logger.info("Starting server on %s...", addr)
	server, err := net.Listen("tcp", addr)
	if err != nil {
		logger.error("Failed to start server: %s", err.Error())
		return false
	}

	logger.info("Started server, waiting for connections")
	c.server = server
	stats.SetServerStartTime()

	return true
}

func (c *Connection) handleConnections(fn handleClientFunc) {
	for {
		client, err := c.server.Accept()

		if err != nil {
			if err_, ok := err.(net.Error); ok && !err_.Temporary() { // server closed
				return
			}

			logger.warning("Failed to accept client: %s", err.Error())
			break
		}

		logger.info("New client connected: %s", client.RemoteAddr().String())

		client.SetReadDeadline(time.Now().Add(time.Duration(config.Read_timeout) * time.Second))
		client.SetWriteDeadline(time.Now().Add(time.Duration(config.Write_timeout) * time.Second))

		stats.IncreaseCurrentConnections()

		clientAbstract := NewClient(client)
		go fn(clientAbstract)
	}
}

func (c *Connection) closeClientConnection(client *Client, reason string) {
	client.connection.Close()

	stats.DecreaseCurrentConnections()

	logger.info("Client disconnect: %s [%s], connected for %s",
		client.connection.RemoteAddr().String(),
		reason,
		client.getConnectionTime())
}
