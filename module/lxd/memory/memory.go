package memory

import (
	"github.com/dppascual/cartobeat/module/lxd"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/common/cfgwarn"
	"github.com/elastic/beats/metricbeat/mb"
	client "github.com/lxc/lxd/client"
)

func init() {
	mb.Registry.MustAddMetricSet("lxd", "memory", New,
		mb.WithHostParser(lxd.HostParser),
		mb.DefaultMetricSet(),
	)
}

// MetricSet stores memory metrics
type MetricSet struct {
	mb.BaseMetricSet
	memoryService *MemoryService
	lxdClient     client.ContainerServer
}

// New creates a new instance of the memory MetricSet
func New(base mb.BaseMetricSet) (mb.MetricSet, error) {
	cfgwarn.Experimental("The lxd memory metricset is experimental.")

	config := lxd.DefaultConfig()
	if err := base.Module().UnpackConfig(&config); err != nil {
		return nil, err
	}

	serverConnection, err := lxd.NewLXDClient(base.HostData().URI, config)
	if err != nil {
		return nil, err
	}

	return &MetricSet{
		BaseMetricSet: base,
		memoryService: &MemoryService{},
		lxdClient:     serverConnection,
	}, nil
}

// Fetch methods implements the data gathering and data conversion to the right
// format. It publishes the event which is then forwarded to the output. In case
// of an error set the Error field of mb.Event or simply call report.Error().
func (m *MetricSet) Fetch(report mb.ReporterV2) {
	report.Event(mb.Event{
		MetricSetFields: common.MapStr{
			"counter": m.counter,
		},
	})
	m.counter++
}
