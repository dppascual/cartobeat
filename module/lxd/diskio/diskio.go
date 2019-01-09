package diskio

import (
	"github.com/dppascual/cartobeat/module/lxd"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/common/cfgwarn"
	"github.com/elastic/beats/metricbeat/mb"
	client "github.com/lxc/lxd/client"
)

func init() {
	mb.Registry.MustAddMetricSet("lxd", "diskio", New,
		mb.WithHostParser(lxd.HostParser),
		mb.DefaultMetricSet(),
	)
}

// MetricSet stores diskio metrics
type MetricSet struct {
	mb.BaseMetricSet
	lxdClient client.ContainerServer
}

// New creates a new instance of the diskio MetricSet
func New(base mb.BaseMetricSet) (mb.MetricSet, error) {
	cfgwarn.Experimental("The lxd diskio metricset is experimental.")

	config := lxd.DefaultConfig()
	if err := base.Module().UnpackConfig(&config); err != nil {
		return nil, err
	}

	serverConnection, err := lxd.NewLXDClient(base.HostData().URI, config, base.Module().Config().Timeout)
	if err != nil {
		return nil, err
	}

	return &MetricSet{
		BaseMetricSet: base,
		lxdClient:     serverConnection,
	}, nil
}

// Fetch creates a list of diskio events for each container.
func (m *MetricSet) Fetch() ([]common.MapStr, error) {
	containerStats, err := lxd.FetchStats(m.lxdClient)
	if err != nil {
		return nil, err
	}

	return eventsMapping(containerStats), nil
}
