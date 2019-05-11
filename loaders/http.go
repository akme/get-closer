package loaders

import (
	"fmt"
	"net/url"
)

func loadFromHTTP(url *url.URL) error {
	var err error
	err = nil
	fmt.Println(url)
	return err
}
