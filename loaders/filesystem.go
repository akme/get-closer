package loaders

import (
	"bufio"
	"os"
)

func loadFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		HostsList = append(HostsList, scanner.Text())
	}
	return scanner.Err()
}
