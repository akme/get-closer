package loaders

import (
	"bufio"
	"fmt"
	"github.com/utahta/go-openuri"
	"io"
)

//LoadHosts loads hosts from file via fs or http request.
func LoadHosts(path string) ([]string, error) {

	o, err := openuri.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer o.Close()

	return linesFromReader(o)
}

func linesFromReader(r io.Reader) ([]string, error) {
	var hostsList []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		hostsList = append(hostsList, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return hostsList, nil
}
