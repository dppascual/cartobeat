package network

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
	events := make([]common.MapStr, len(containerStat.State.Network))

	for key, value := range containerStat.State.Network {
		event := common.MapStr{
			mb.ModuleDataKey: common.MapStr{
				"container": common.MapStr{
					"name": containerStat.Container.Name,
				},
			},
			"interface": key,
			"inbound" : common.MapStr{
				"bytes": value.Counters.BytesReceived,
				"packets": value.Counters.PacketsReceived,
			},
			"outbound": common.MapStr{
				"bytes": value.Counters.BytesSent,
				"packets": value.Counters.PacketsSent,
			},
		}
		events = append(events, event)
	}
	return events
}
