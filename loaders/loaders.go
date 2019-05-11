package loaders

import (
	"fmt"
	"net/url"
)

type Hosts []string

var HostsList Hosts

func LoadHosts(path string) error {
	// parse connection string
	u, err := url.Parse(path) //(viper.GetString("backend"))
	if err != nil {
		return err
	}

	// set backend based on connection string's scheme
	switch u.Scheme {
	case "file", "":
		err = loadFromFile(u.Path) // &Filesystem{Path: u.Path}
	case "http", "https":
		err = loadFromHTTP(u)
	default:
		return fmt.Errorf(`
Unrecognized scheme '%s'. You can visit https://github.com/nanopack/hoarder and
submit a pull request adding the scheme or you can submit an issue requesting its
addition.
`, u.Scheme)
	}

	// initialize the driver
	return err
}
