package metrics

import (
	"github.com/pritunl/pritunl-monitoring/database"
	"github.com/pritunl/pritunl-monitoring/hosts"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	cpuUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "pritunl_cpu_usage",
		Help: "Current CPU usage of Pritunl process",
	})
	memUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "pritunl_mem_usage",
		Help: "Current memory usage of Pritunl process",
	})
	deviceCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "pritunl_connected_devices",
		Help: "Current number of devices connected to Pritunl node",
	})
	serverCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "pritunl_running_servers",
		Help: "Current number of servers running to Pritunl node",
	})
	threadCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "pritunl_thread_count",
		Help: "Current number of threads in Pritunl process",
	})
)

func init() {
	prometheus.MustRegister(cpuUsage)
	prometheus.MustRegister(memUsage)
	prometheus.MustRegister(deviceCount)
	prometheus.MustRegister(serverCount)
	prometheus.MustRegister(threadCount)
}

func Update() (err error) {
	db := database.GetDatabase()
	defer db.Close()

	hsts, err := hosts.GetHosts(db)
	if err != nil {
		return
	}

	for _, host := range hsts {
		if host.Status == "online" {
			cpuUsage.Set(host.CpuUsage)
			memUsage.Set(host.MemUsage)
			serverCount.Set(float64(host.ServerCount))
			deviceCount.Set(float64(host.DeviceCount))
			threadCount.Set(float64(host.ThreadCount))
		} else {
			cpuUsage.Set(0.0)
			memUsage.Set(0.0)
			serverCount.Set(0.0)
			deviceCount.Set(0.0)
			threadCount.Set(0.0)
		}
	}

	return
}
