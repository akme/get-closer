package loaders

import (
	"bufio"
	"fmt"
	"github.com/utahta/go-openuri"
	"io"
)

type Hosts []string

var HostsList Hosts

func LoadHosts(path string) error {

	o, err := openuri.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer o.Close()

	return linesFromReader(o)
}

func linesFromReader(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		HostsList = append(HostsList, scanner.Text())
	}
	var err error
	if err := scanner.Err(); err != nil {
		return err
	}

	return err
}
