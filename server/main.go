// socket-server project main.go
package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

// TODO: should be share in a lib
const BUFFER_SIZE = 1024

type Client struct {
	name       string
	connection net.Conn
	ch         chan string
}

var clientMaps = SafeMap[string, Client]{v: make(map[string]*Client)}

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
	client := Client{connection: connection}

	for {
		command, body, err := ReadCommand(connection)
		if err != nil {
			fmt.Println("Error while reading a command:", err.Error())
			panic("Error while receiving a command")
		}

		switch command {
		case COMMAND_SET_NAME:
			{
				fmt.Printf("\t- [%s]: %s\n", COMMAND_SET_NAME, body[0])
				client.name = body[0]
				client.ch = make(chan string)
				clientMaps.Add(client.name, &client)
			}
		default:
			panic("Unknown command received:" + command)
		}
	}

	buffer := make([]byte, BUFFER_SIZE)
	mLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	fmt.Println("Received: ", string(buffer[:mLen]))
	_, err = connection.Write([]byte("Thanks! Got your message:" + string(buffer[:mLen])))
	if err != nil {
		panic("Error while sending a message")
	}
	connection.Close()
}

// TODO: should be in a shared module
func getEnvOrPanic(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("%s must be defined", key))
	}
	return val
}

// TODO: should be in a shared module
const COMMAND_SET_NAME = "SET_NAME"

// returns the command name, the associated values and an error.
// TODO: should be in a shared module
func ReadCommand(connection net.Conn) (string, []string, error) {
	buffer := make([]byte, BUFFER_SIZE)

	readLine := func() (string, error) {
		mLen, err := connection.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			return "", err
		}
		// -- ack
		_, err = connection.Write(make([]byte, 1))
		if err != nil {
			panic(err)
		}

		return string(buffer[:mLen]), nil
	}

	// - number of lines
	linesStr, err := readLine()
	if err != nil {
		return "", nil, err
	}
	lines, err := strconv.Atoi(linesStr)
	if err != nil {
		return "", nil, err
	}
	// - command
	command, err := readLine()
	if err != nil {
		return "", nil, err
	}
	// - body
	body := make([]string, lines)
	for i := 0; i < lines; i++ {
		bodyLine, err := readLine()
		if err != nil {
			return "", nil, err
		}
		body[i] = bodyLine
	}

	return command, body, nil
}

// TODO: it should exist a lib doing this better
type SafeMap[K comparable, V interface{}] struct {
	mu sync.Mutex
	v  map[K]*V
}

func (safeMap *SafeMap[K, V]) Add(key K, value *V) *V {
	safeMap.mu.Lock()
	defer safeMap.mu.Unlock()
	safeMap.v[key] = value
	return value
}

func (safeMap *SafeMap[K, V]) Get(key K) *V {
	safeMap.mu.Lock()
	defer safeMap.mu.Unlock()
	return safeMap.v[key]
}

func (safeMap *SafeMap[K, V]) Del(key K) {
	safeMap.mu.Lock()
	defer safeMap.mu.Unlock()
	delete(safeMap.v, key)
}
