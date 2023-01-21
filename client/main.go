package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
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
	// buffer := make([]byte, BUFFER_SIZE)

	// send some data
	// - name
	err = SendSetName(connection, "Fabien")
	if err != nil {
		panic(err)
	}
	fmt.Println("Name sent")
	// - message
	err = SendMessage(connection, "Coucou Amrit!")
	if err != nil {
		panic(err)
	}
	fmt.Println("Message sent")

	// mLen, err := connection.Read()
	// if err != nil {
	// 	fmt.Println("Error reading:", err.Error())
	// }
	// fmt.Println("Received: ", string(buffer[:mLen]))
}

func SendMessage(connection net.Conn, message string) error {
	err := SendWithAck(connection, []byte("2"))
	if err != nil {
		return err
	}
	err = SendWithAck(connection, []byte(COMMAND_SEND_MESSAGE))
	if err != nil {
		return err
	}
	err = SendWithAck(connection, []byte(message))
	if err != nil {
		return err
	}
	err = SendWithAck(connection, []byte(strconv.FormatInt(time.Now().UnixMilli(), 10)))
	if err != nil {
		return err
	}
	return nil
}

func SendSetName(connection net.Conn, name string) error {
	err := SendWithAck(connection, []byte("1"))
	if err != nil {
		return err
	}
	err = SendWithAck(connection, []byte(COMMAND_SET_NAME))
	if err != nil {
		return err
	}
	err = SendWithAck(connection, []byte(name))
	if err != nil {
		return err
	}
	return nil
}

func SendWithAck(connection net.Conn, bytes []byte) error {
	_, err := connection.Write(bytes)
	if err != nil {
		return err
	}
	_, err = connection.Read(make([]byte, 1))
	if err != nil {
		return err
	}
	return nil
}

// TODO: should be in a shared module
func getEnvOrPanic(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("%s must be defined", key))
	}
	return val
}
