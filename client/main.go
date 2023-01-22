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

	// establish connection
	connection, err := net.Dial(serverType, serverHost+":"+serverPort)
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	// goroutine to receive messages from server
	go ListenServer(connection)

	// send some data
	// - name
	name, err := utils.ReadFromStdin("Enter your name: ")
	if err != nil {
		panic(err)
	}
	err = com.SendSetName(connection, name)
	if err != nil {
		panic(err)
	}
	fmt.Println("Name sent")
	// - messages
	// 		/q to quit
	for {
		dest, err := utils.ReadFromStdin("Enter your message dest: ")
		if err != nil {
			panic(err)
		}
		if dest == "/q" {
			fmt.Println("Quitting...")
			return
		}
		message, err := utils.ReadFromStdin("Enter your message: ")
		if err != nil {
			panic(err)
		}
		if message == "/q" {
			fmt.Println("Quitting...")
			return
		}
		err = com.SendMessage(connection, dest, message)
		if err != nil {
			panic(err)
		}
		fmt.Println("Message sent")
	}
}

func ListenServer(connection net.Conn) {
	for {
		command, body, err := com.ReadCommand(connection)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Server has disconnected\n")
				return
			}
			fmt.Println("Error while reading a command:", err.Error())
			panic("Error while receiving a command")
		}

		switch command {
		case com.COMMAND_RECEIVE_MESSAGE:
			{
				fmt.Printf("\t- [%s]: %s <- %s @%s\n", com.COMMAND_RECEIVE_MESSAGE, body[1], body[0], body[2])
			}
		default:
			panic("Unknown command received:" + command)
		}
	}
}

