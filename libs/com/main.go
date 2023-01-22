package com

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

const COMMAND_SET_NAME = "SET_NAME"
const COMMAND_SEND_MESSAGE = "SEND_MESSAGE"
const COMMAND_RECEIVE_MESSAGE = "RECEIVE_MESSAGE"

const BUFFER_SIZE = 1024

func SendMessage(connection net.Conn, dest string, message string) error {
	return SendCommand(connection, []string{COMMAND_SEND_MESSAGE, dest, message, strconv.FormatInt(time.Now().UnixMilli(), 10)})
}

func SendSetName(connection net.Conn, name string) error {
	return SendCommand(connection, []string{COMMAND_SET_NAME, name})
}

func SendReceiveMessage(connection net.Conn, name string, message string, timestamp string) error {
	return SendCommand(connection, []string{COMMAND_RECEIVE_MESSAGE, name, message, timestamp})
}

// returns the command name, the associated values and an error.
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

func SendCommand(connection net.Conn, parts []string) error {
	_, err := connection.Write([]byte(strings.Join(parts, "%part%\n") + "%end%\n"))

	return err
}
