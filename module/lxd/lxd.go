package lxd

import (
	"errors"
	"net/url"
	"sync"
	"time"

	"github.com/elastic/beats/metricbeat/mb"
	"github.com/elastic/beats/metricbeat/mb/parse"
	client "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
)

// LXD API version
const lxdAPIVersion = "1.0"

// HostParser returns a function that parses a host value from the configuration
// and returns a HostData object.
var HostParser = parse.URLHostParserBuilder{DefaultScheme: "tcp"}.Build()

func init() {
	// Register the ModuleFactory function for the "lxd" module.
	if err := mb.Registry.AddModule("lxd", NewModule); err != nil {
		panic(err)
	}
}

// NewModule is used to initialize the "lxd" module.
func NewModule(base mb.BaseModule) (mb.Module, error) {
	// Validate that at least one host has been specified.
	config := struct {
		Hosts []string `config:"hosts"    validate:"nonzero,required"`
	}{}
	if err := base.UnpackConfig(&config); err != nil {
		return nil, err
	}

	return &base, nil
}

// ContainerStats contains the stats of all containers up
type ContainerStats struct {
	Container *api.Container
	State     *api.ContainerState
}

// NewLXDClient initializes and returns a new connection to a LXD server.
func NewLXDClient(endpoint string, config Config, timeout time.Duration) (client.ContainerServer, error) {
	var c client.ContainerServer

	urlRaw, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	switch urlRaw.Scheme {
	case "unix":
		c, err = client.ConnectLXDUnix(urlRaw.Path, nil)
		if err != nil {
			return nil, err
		}
	case "https":
		connArgs, err := initHTTPSConn(config, timeout)
		if err != nil {
			return nil, err
		}
		c, err = client.ConnectLXD(urlRaw.Path, connArgs)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("The connection is only allowed over a local unix socket or HTTPs")
	}

	return c, nil
}

// FetchStats returns a list of active containers with all related stats inside
func FetchStats(serverConnection client.ContainerServer) ([]ContainerStats, error) {

	containersList, err := FetchInfo(serverConnection)
	if err != nil {
		return nil, err
	}

	// Declare a new slice to store active containers
	containersActive := make([]ContainerStats, 0, len(containersList))

	for _, stat := range containersList {
		if stat.Container.IsActive() && stat.State != nil {
			containersActive = append(containersActive, stat)
		}
	}

	return containersActive, nil
}

// FetchInfo returns a list of all containers with all related stats inside
func FetchInfo(serverConnection client.ContainerServer) ([]ContainerStats, error) {

	containers, err := serverConnection.GetContainers()
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup

	containersList := make([]ContainerStats, 0, len(containers))
	statsQueue := make(chan ContainerStats)
	wg.Add(len(containers))

	for _, container := range containers {
		go func(container api.Container) {
			defer wg.Done()
			statsQueue <- getContainerStats(serverConnection, &container)
		}(container)
	}

	go func() {
		wg.Wait()
		close(statsQueue)
	}()

	for stat := range statsQueue {
		containersList = append(containersList, stat)
	}

	return containersList, nil
}
