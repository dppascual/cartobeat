package memory

import (
	"github.com/dppascual/cartobeat/module/lxd"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/metricbeat/mb"
)

func eventsMapping(containerStats []lxd.ContainerStats) []common.MapStr {
	formattedEvents := []common.MapStr{}
	for _, containerStat := range containerStats {
		formattedEvents = append(formattedEvents, eventMapping(containerStat))
	}
	return formattedEvents
}

func eventMapping(containerStat lxd.ContainerStats) common.MapStr {
	event := common.MapStr{
		mb.ModuleDataKey: common.MapStr{
			"container": containerStat.Container.Name,
		},
		"usage": common.MapStr{
			"total": containerStat.State.Memory.Usage,
			"peak":  containerStat.State.Memory.UsagePeak,
		},
		"swap": common.MapStr{
			"total": containerStat.State.Memory.SwapUsage,
			"peak":  containerStat.State.Memory.SwapUsagePeak,
		},
	}
	return event
}
