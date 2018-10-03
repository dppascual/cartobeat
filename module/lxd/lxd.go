package lxd

import (
	"errors"
	"net/url"

	"github.com/elastic/beats/metricbeat/mb/parse"
	client "github.com/lxc/lxd/client"
)

// LXD API version
const lxdAPIVersion = "1.0"

// HostParser returns a function that parses a host value from the configuration
// and returns a HostData object.
var HostParser = parse.URLHostParserBuilder{DefaultScheme: "tcp"}.Build()

// NewLXDClient initializes and returns a new connection to a LXD server
func NewLXDClient(endpoint string, config Config) (client.ContainerServer, error) {
	var c client.ContainerServer

	urlRaw, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	switch urlRaw.Scheme {
	case "unix":
		c, err = client.ConnectLXDUnix(endpoint, nil)
		if err != nil {
			return nil, err
		}
	case "https":
		connArgs, err := initHTTPSConn(config)
		if err != nil {
			return nil, err
		}
		c, err = client.ConnectLXD(endpoint, connArgs)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("The connection is only allowed over a local unix socket or HTTPs")
	}

	return c, nil
}
