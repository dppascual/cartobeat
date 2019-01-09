package diskio

import (
	"github.com/dppascual/cartobeat/module/lxd"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/metricbeat/mb"
)

func eventsMapping(containerStats []lxd.ContainerStats) []common.MapStr {
	formattedEvents := []common.MapStr{}
	for _, containerStat := range containerStats {
		for _, eventDisk := range eventMapping(containerStat) {
			formattedEvents = append(formattedEvents, eventDisk)
		}
	}
	return formattedEvents
}

func eventMapping(containerStat lxd.ContainerStats) []common.MapStr {
	events := make([]common.MapStr, len(containerStat.State.Disk))

	for key, value := range containerStat.State.Disk {
		event := common.MapStr{
			mb.ModuleDataKey: common.MapStr{
				"container": common.MapStr{
					"name": containerStat.Container.Name,
				},
			},
			"name": key,
			"usage": value.Usage,
		}

		events = append(events, event)
	}

	return events
}
