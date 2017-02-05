package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/edwinhoksberg/tamarin/command"
	"github.com/edwinhoksberg/tamarin/message"
)

const APP_NAME string = "tamarin"
const APP_VERSION string = "1.0.0"

var connection *Connection
var logger *Logger
var config *Config
var stats *Stats

func main() {
	configFile := flag.String("config", "./config.json", "The tamarin config file")
	flag.Parse()

	catchSignals()

	config = new(Config)
	if !config.loadConfig(*configFile) {
		os.Exit(1)
	}

	logger = NewLogger()
	stats = NewStats()

	logger.info("Starting with PID %d...", os.Getpid())

	connection = new(Connection)

	connection.start(config.Listen, config.Port)
	connection.handleConnections(connectionHandler)
}

func connectionHandler(client *Client) {
	buffer := bufio.NewReader(client.connection)
	defer client.connection.Close()

	for {
		if client.isNewConnection {
			// if new connection: 200 news.php.net - colobus 2.1 ready - (posting ok).
			welcomeMsg := fmt.Sprintf("200 %s - %s %s - posting allowed\r\n", config.Hostname, APP_NAME, APP_VERSION)

			client.connection.Write([]byte(welcomeMsg))
			client.isNewConnection = false
		}

		// Read until we encouter a newline
		line, err := buffer.ReadBytes('\n')
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				connection.closeClientConnection(client, "timed out")
				return
			}

			if err == io.EOF {
				connection.closeClientConnection(client, "connection closed")
				return
			}

			logger.warning("Recieved invalid client data: %s", err.Error())
			continue
		}

		logger.debug("Recieved message: %s", line)

		// Parse the client message in multiple parameters
		recieved_message := strings.Fields(string(line[:]))
		if len(recieved_message) == 0 {
			logger.info("Got empty client message, ignoring...")
			continue
		}

		// Make an request object with the parsed message
		logger.debug("Generating request object...")
		request := message.NewRequest(strings.ToLower(recieved_message[0]), append(recieved_message[:0], recieved_message[:1]...))

		// Find the associated command for the request
		logger.debug("Generating request message...")
		generated_command := command.New(request).GenerateResponse()

		// Generate an response for the client
		response := generated_command.ToString()

		// Send the response
		logger.debug("Sending message: %s", strings.Trim(response, "\r\n"))
		_, err = client.connection.Write([]byte(response))
		if err != nil {
			logger.error("Could not send message: %s", err.Error())
		} else {
			logger.debug("Message succesfully sent")
		}

		// If it was an QUIT command, close the connection
		if request.GetCommand() == "quit" {
			connection.closeClientConnection(client, "quit")

			return
		}

		client.connection.SetReadDeadline(time.Now().Add(time.Duration(config.Read_timeout) * time.Second))
		client.connection.SetWriteDeadline(time.Now().Add(time.Duration(config.Write_timeout) * time.Second))
	}
}

func catchSignals() {
	signal_channel := make(chan os.Signal, 1)

	signal.Notify(signal_channel, os.Interrupt, syscall.SIGTERM, syscall.SIGUSR1)
	go func() {
		for {
			recv_signal := <-signal_channel

			switch recv_signal {
			case syscall.SIGINT:
				fallthrough
			case syscall.SIGTERM:
				logger.info("Recieved program interrupt, shutting down...")
				connection.server.Close()

				os.Exit(0)
			case syscall.SIGUSR1:
				stats.PrintStats()
			}
		}
	}()
}
