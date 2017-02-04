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

	//connections := make(chan *Client)
	//
	//go func() {
	//	for {
	//		client, err := c.server.Accept()
	//		if err != nil {
	//			if err_, ok := err.(net.Error); ok && !err_.Temporary() {
	//				// srever closed
	//				close(connections)
	//
	//				return
	//			}
	//
	//			logger.warning("Failed to accept client: %s", err.Error())
	//			break
	//		}
	//		logger.info("New client connected: %s", client.RemoteAddr().String())
	//
	//		client.SetReadDeadline(time.Now().Add(time.Duration(config.Read_timeout) * time.Second))
	//		client.SetWriteDeadline(time.Now().Add(time.Duration(config.Write_timeout) * time.Second))
	//
	//		atomic.AddUint32(&c.connectedClients, 1)
	//
	//		currentCount := atomic.LoadUint32(&c.connectedClients)
	//		maxCount := atomic.LoadUint32(&c.connectedClientsMax)
	//
	//		if currentCount > maxCount {
	//			atomic.AddUint32(&c.connectedClientsMax, currentCount - maxCount)
	//		}
	//
	//	// @todo only call this when daemon asks, remove this later
	//		logger.stats("Server uptime: %s", time.Since(connection.serverStartTime))
	//		logger.stats("Current connections: %s", strconv.Itoa(int(atomic.LoadUint32(&c.connectedClients))))
	//		logger.stats("Maximum concurrent connections: %s", strconv.Itoa(int(atomic.LoadUint32(&c.connectedClientsMax))))
	//
	//		var mem runtime.MemStats
	//		runtime.ReadMemStats(&mem)
	//		logger.stats("Current memory allocated: %d kB", mem.Alloc / 1024)
	//
	//		clientAbstract := fromConnection(client)
	//		logger.debug("Created new client object")
	//
	//		connections <- clientAbstract
	//	}
	//}()
	//
	//for {
	//	go fn(<-connections)
	//}
}

func (c *Connection) closeClientConnection(client *Client, reason string) {
	client.connection.Close()

	stats.DecreaseCurrentConnections()

	logger.info("Client disconnect: %s [%s], connected for %s",
		client.connection.RemoteAddr().String(),
		reason,
		client.getConnectionTime())
}
