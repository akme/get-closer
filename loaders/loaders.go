package loaders

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
)

type Hosts []string

var HostsList Hosts

func LoadHosts(path string) error {
	// parse connection string
	u, err := url.Parse(path)
	if err != nil {
		return err
	}
	switch u.Scheme {
	case "file", "":
		err = loadFromFile(u.Path) // &Filesystem{Path: u.Path}
	case "http", "https":
		err = loadFromHTTP(path)
	default:
		return fmt.Errorf(`
Unrecognized scheme '%s'. You can visit https://github.com/nanopack/hoarder and
submit a pull request adding the scheme or you can submit an issue requesting its
addition.
`, u.Scheme)
	}

	return err
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
