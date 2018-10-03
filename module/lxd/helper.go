package lxd

import (
	"errors"
	"net/http"
	"time"

	client "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared"
	"github.com/lxc/lxd/shared/api"
)

// initHTTPSConn charges TLS config into a structure that is used by lxd client.
func initHTTPSConn(config Config, timeout time.Duration) (*client.ConnectionArgs, error) {
	var options client.ConnectionArgs

	if config.TLS.IsEnabled() {
		options = client.ConnectionArgs{
			TLSCA:         config.TLS.CA,
			TLSClientCert: config.TLS.ClientCert,
			TLSClientKey:  config.TLS.ClientKey,
		}

		tlsc, err := shared.GetTLSConfig(options.TLSClientCert, options.TLSClientKey, options.TLSCA, nil)
		if err != nil {
			return &client.ConnectionArgs{}, err
		}

		options.HTTPClient = &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				TLSClientConfig: tlsc,
			},
		}
	} else {
		return &client.ConnectionArgs{}, errors.New("A TLS config must be provided to connect to a remote LXD daemon over HTTPs")
	}

	return &options, nil
}

// getContainerStats loads stats for the given container.
func getContainerStats(serverConnection client.ContainerServer, container *api.Container) ContainerStats {

	event := ContainerStats{
		Container: container,
	}

	containerState, _, err := serverConnection.GetContainerState(container.Name)
	if err != nil {
		return event
	}

	event.State = containerState
	return event
}
