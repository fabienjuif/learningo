package main

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/joho/godotenv"

	"github.com/fabienjuif/learningo/libs/com"
	"github.com/fabienjuif/learningo/libs/env"
	"github.com/fabienjuif/learningo/libs/utils"
)

type Client struct {
	name       string
	connection net.Conn
}

var clientMaps = utils.NewSafeMap(make(map[string]*Client))

func main() {
	// TODO: all this init stuff should be in a shared module
	envName := os.Getenv("ENV")
	err := godotenv.Load(envName + ".env")
	if err != nil {
		panic("Error loading .env file")
	}
	serverHost := env.GetEnvOrPanic("SERVER_HOST")
	serverPort := env.GetEnvOrPanic("SERVER_PORT")
	serverType := env.GetEnvOrPanic("SERVER_TYPE")

	fmt.Println("Server Running...")
	server, err := net.Listen(serverType, serverHost+":"+serverPort)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer server.Close()
	fmt.Println("Listening on " + serverHost + ":" + serverPort)
	fmt.Println("Waiting for client...")
	for {
		connection, err := server.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("client connected")
		go processClient(connection)
	}
}

func processClient(connection net.Conn) {
	defer connection.Close()
	var client *Client

	for {
		command, body, err := com.ReadCommand(connection)
		if err != nil {
			if err == io.EOF {
				if client != nil {
					clientMaps.Del(client.name)
					fmt.Printf("Client has disconnected: %s\n", client.name)
					return
				}
			}
			fmt.Println("Error while reading a command:", err.Error())
			panic("Error while receiving a command")
		}

		switch command {
		case com.COMMAND_SET_NAME:
			{
				fmt.Printf("\t- [%s]: %s\n", com.COMMAND_SET_NAME, body[0])
				name := body[0]
				client = &Client{name, connection}
				_, err := clientMaps.Add(name, client)
				if err != nil {
					return
				}

			}
		case com.COMMAND_SEND_MESSAGE:
			{
				fmt.Printf("\t- [%s]: %s <- %s @%s\n", com.COMMAND_SEND_MESSAGE, body[0], body[1], body[2])
				go func() {
					dest := clientMaps.Get(body[0])
					if dest == nil {
						fmt.Printf("Client not found: %s\n", body[0])
						return
					}
					err := com.SendReceiveMessage(dest.connection, client.name, body[1], body[2])
					if err != nil {
						fmt.Println("Error while sending message", err.Error())
					}
				}()
			}
		default:
			panic("Unknown command received:" + command)
		}
	}
}
