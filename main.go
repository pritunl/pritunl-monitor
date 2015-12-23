package main

import (
	"fmt"
	"github.com/pritunl/pritunl-monitor/datadog"
	"github.com/pritunl/pritunl-monitor/prometheus"
	"os"
)

func main() {
	switch os.Getenv("MODE") {
	case "prometheus":
		err := prometheus.Start()
		if err != nil {
			panic(err)
		}
	case "datadog":
		err := datadog.Start()
		if err != nil {
			panic(err)
		}
	default:
		panic(fmt.Sprintf("unknown mode %s", os.Getenv("MODE")))
	}
}
