package loaders

import (
	"net/http"
)

func loadFromHTTP(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = linesFromReader(resp.Body)

	return err
}
