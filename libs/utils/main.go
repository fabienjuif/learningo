package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReadFromStdin(description string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(description)
	data, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(data, "\n"), nil
}