package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// TODO: should be share in a lib
const BUFFER_SIZE = 1024

// TODO: should be in a shared module
const COMMAND_SET_NAME = "SET_NAME"
const COMMAND_SEND_MESSAGE = "SEND_MESSAGE"

func main() {
	// TODO: all this init stuff should be in a shared module
	envName := os.Getenv("ENV")
	err := godotenv.Load(envName + ".env")
	if err != nil {
		panic("Error loading .env file")
	}
	serverHost := getEnvOrPanic("SERVER_HOST")
	serverPort := getEnvOrPanic("SERVER_PORT")
	serverType := getEnvOrPanic("SERVER_TYPE")

	// establish connection
	connection, err := net.Dial(serverType, serverHost+":"+serverPort)
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	// send some data
	// - name
	name, err := readFromStdin("Enter your name: ")
	if err != nil {
		panic(err)
	}
	err = SendSetName(connection, name)
	if err != nil {
		panic(err)
	}
	fmt.Println("Name sent")
	// - messages
	// 		/q to quit
	for {
		dest, err := readFromStdin("Enter your message dest: ")
		if err != nil {
			panic(err)
		}
		if dest == "/q" {
			fmt.Println("Quitting...")
			return
		}
		message, err := readFromStdin("Enter your message: ")
		if err != nil {
			panic(err)
		}
		if message == "/q" {
			fmt.Println("Quitting...")
			return
		}
		err = SendMessage(connection, dest, message)
		if err != nil {
			panic(err)
		}
		fmt.Println("Message sent")
	}
}

func SendMessage(connection net.Conn, dest string, message string) error {
	return Send(connection, []string{COMMAND_SEND_MESSAGE, dest, message, strconv.FormatInt(time.Now().UnixMilli(), 10)})
}

func SendSetName(connection net.Conn, name string) error {
	return Send(connection, []string{COMMAND_SET_NAME, name})
}

func Send(connection net.Conn, parts []string) error {
	_, err := connection.Write([]byte(strings.Join(parts, "%part%\n") + "%end%\n"))

	return err
}

// TODO: should be in a shared module
func getEnvOrPanic(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("%s must be defined", key))
	}
	return val
}

func readFromStdin(description string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(description)
	data, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(data, "\n"), nil
}
