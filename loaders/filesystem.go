package loaders

import (
	"os"
)

func loadFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	err = linesFromReader(file)
	return err
}
