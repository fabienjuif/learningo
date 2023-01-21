package main

import (
	"fmt"
	"net"
	"os"

	"github.com/joho/godotenv"
)

// TODO: should be share in a lib
const BUFFER_SIZE = 1024

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

	//establish connection
	connection, err := net.Dial(serverType, serverHost+":"+serverPort)
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	buffer := make([]byte, BUFFER_SIZE)
	ackBuffer := make([]byte, 1)

	///send some data
	_, err = connection.Write([]byte("1"))
	if err != nil {
		panic(err)
	}
	_, err = connection.Read(ackBuffer)
	if err != nil {
		panic(err)
	}
	_, err = connection.Write([]byte("SET_NAME"))
	if err != nil {
		panic(err)
	}
	_, err = connection.Read(ackBuffer)
	if err != nil {
		panic(err)
	}
	_, err = connection.Write([]byte("Fabien"))
	if err != nil {
		panic(err)
	}

	mLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	fmt.Println("Received: ", string(buffer[:mLen]))
}

// TODO: should be in a shared module
func getEnvOrPanic(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("%s must be defined", key))
	}
	return val
}
