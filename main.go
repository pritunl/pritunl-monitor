package main

import (
	"github.com/pritunl/pritunl-monitoring/prometheus"
)

func main() {
	err := prometheus.Start()
	if err != nil {
		panic(err)
	}
}
