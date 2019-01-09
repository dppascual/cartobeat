package main

import (
	"os"

	"github.com/elastic/beats/libbeat/cmd/instance"
	"github.com/elastic/beats/metricbeat/beater"

	// Comment out the following line to exclude all official metricbeat modules and metricsets
	_ "github.com/elastic/beats/metricbeat/include"

	// Make sure all your modules and metricsets are linked in this file
	_ "github.com/dppascual/cartobeat/include"
)

var Name = "cartobeat"
var IndexPrefix = "cartobeat"
var Version = "0.2.0"

func main() {
	if err := instance.Run(Name, IndexPrefix, Version, beater.DefaultCreator()); err != nil {
		os.Exit(1)
	}
}
