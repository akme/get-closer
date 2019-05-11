package loaders

import (
	"bufio"
	"fmt"
	"github.com/utahta/go-openuri"
	"io"
	"net/url"
	"strings"
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
		s := scanner.Text()
		// if blank line or comment, skip
		if s == "" || strings.HasPrefix(s, "#") {
			continue
		}

		if strings.Contains(s, "http://") || strings.Contains(s, "https://") {
			host, err := url.Parse(s)
			if err != nil {
				continue
			}
			hostsList = append(hostsList, host.Hostname())
		} else {
			hostsList = append(hostsList, s)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return hostsList, nil
}
