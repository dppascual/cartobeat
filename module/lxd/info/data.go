package info

import (
	"github.com/dppascual/cartobeat/module/lxd"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/metricbeat/mb"
	"github.com/lxc/lxd/shared/api"
)

func eventsMapping(containerStats []lxd.ContainerStats) []common.MapStr {
	containers := map[api.StatusCode]int{
		102: 0,
		103: 0,
	}

	for _, containerStat := range containerStats {
		if _, ok := containers[containerStat.Container.StatusCode]; ok {
			containers[containerStat.Container.StatusCode] += 1
		}
	}

	return []common.MapStr{
		common.MapStr{
			mb.ModuleDataKey: common.MapStr{
				"containers": common.MapStr{
					"running": containers[103],
					"stopped": containers[102],
					"total":   containers[102] + containers[103],
				},
			},
		},
	}
}
