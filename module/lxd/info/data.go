package info

import (
	"github.com/dppascual/cartobeat/module/lxd"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/metricbeat/mb"
)

func eventsMapping(containerStats []lxd.ContainerStats) []common.MapStr {
	var containersRunning int
	var containersStopped int

	for _, containerStat := range containerStats {
		if containerStat.State.StatusCode == 103 {
			containersRunning += 1
		} else if containerStat.State.StatusCode == 102 {
			containersStopped += 1
		}
	}

	return []common.MapStr{
		common.MapStr{
			mb.ModuleDataKey: common.MapStr{
				"containers": common.MapStr{
					"running": containersRunning,
					"stopped": containersStopped,
					"total":   containersRunning + containersStopped,
				},
			},
		},
	}
}
