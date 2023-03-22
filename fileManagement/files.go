package fileManagement

import (
	"bufio"
	"fmt"
	"os"
)

func WriteLine(data string, fileNameWithExtension string) error {
	file, err := os.Create(fileNameWithExtension)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
	return err
}

func ReadLine(fileNameWithExtension string) (string, error) {
	file, err := os.Open(fileNameWithExtension)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var line string
	for scanner.Scan() {
		line = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	return line, err
}
