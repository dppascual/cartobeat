package lxd

import (
	"errors"
	"net/http"

	client "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared"
)

func initHTTPSConn(config Config) (*client.ConnectionArgs, error) {
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
			Transport: &http.Transport{
				TLSClientConfig: tlsc,
			},
		}
	} else {
		return &client.ConnectionArgs{}, errors.New("A TLS config must be provided to connect to a remote LXD daemon over HTTPs")
	}

	return &options, nil
}
