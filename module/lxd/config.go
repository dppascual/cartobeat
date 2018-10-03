package lxd

// Config stores the LXD config
type Config struct {
	TLS *TLSConfig `config:"ssl"`
}

// TLSConfig is used to store the LXD TLS config
type TLSConfig struct {
	Enabled    bool   `config:"enabled"`
	CA         string `config:"certificate_authority"`
	ClientCert string `config:"certificate"`
	ClientKey  string `config:"key"`
}

// DefaultConfig returns default module config
func DefaultConfig() Config {
	return Config{}
}

// IsEnabled returns if TLS is enabled
func (c *TLSConfig) IsEnabled() bool {
	return c != nil && c.Enabled
}
