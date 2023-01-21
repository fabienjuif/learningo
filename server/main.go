// socket-server project main.go
package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

// TODO: should be share in a lib
const BUFFER_SIZE = 1024
const PACKET_BUFFER_SIZE = 2048

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
	defer connection.Close()
	var client *Client

	for {
		command, body, err := ReadCommand(connection)
		if err != nil {
			if err == io.EOF {
				if client != nil {
					clientMaps.Del(client.name)
					fmt.Printf("Client has disconnected: %s\n", client.name)
					close(client.ch)
					return
				}
			}
			fmt.Println("Error while reading a command:", err.Error())
			panic("Error while receiving a command")
		}

		switch command {
		case COMMAND_SET_NAME:
			{
				fmt.Printf("\t- [%s]: %s\n", COMMAND_SET_NAME, body[0])
				name := body[0]
				client = &Client{name, connection, make(chan string)}
				_, err := clientMaps.Add(name, client)
				if err != nil {
					return
				}

			}
		case COMMAND_SEND_MESSAGE:
			{
				fmt.Printf("\t- [%s]: %s -> %s @%s\n", COMMAND_SEND_MESSAGE, body[0], body[1], body[2])
				// FIXME: remove it this is just a test
				err := Send(client.connection, []string{COMMAND_RECEIVE_MESSAGE, "Amrit", body[1], body[2]})
				if err != nil {
					panic(err.Error())
				}
			}
		default:
			panic("Unknown command received:" + command)
		}
	}
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
const COMMAND_SEND_MESSAGE = "SEND_MESSAGE"
const COMMAND_RECEIVE_MESSAGE = "RECEIVE_MESSAGE"

// returns the command name, the associated values and an error.
// TODO: should be in a shared module
func ReadCommand(connection net.Conn) (string, []string, error) {
	buffer := make([]byte, BUFFER_SIZE)
	message := ""

	for !strings.HasSuffix(message, "%end%\n") {
		mLen, err := connection.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			return "", nil, err
		}
		message += string(buffer[:mLen])
	}
	split := strings.Split(strings.TrimSuffix(message, "%end%\n"), "%part%\n")
	return split[0], split[1:], nil
}

// TODO: should be in a shared module
func Send(connection net.Conn, parts []string) error {
	_, err := connection.Write([]byte(strings.Join(parts, "%part%\n") + "%end%\n"))

	return err
}

// TODO: it should exist a lib doing this better
type SafeMap[K comparable, V interface{}] struct {
	mu sync.Mutex
	v  map[K]*V
}

func (safeMap *SafeMap[K, V]) Add(key K, value *V) (*V, error) {
	safeMap.mu.Lock()
	defer safeMap.mu.Unlock()
	_, exists := safeMap.v[key]
	if exists {
		return nil, &SafeMapAlreadyExistsError{}
	}
	safeMap.v[key] = value
	return value, nil
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

type SafeMapAlreadyExistsError struct{}

func (m *SafeMapAlreadyExistsError) Error() string {
	return "Key already exists"
}
