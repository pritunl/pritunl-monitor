package datadog

import (
	"bytes"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/dropbox/godropbox/errors"
	"github.com/dustin/gojson"
	"github.com/pritunl/pritunl-monitoring/database"
	"github.com/pritunl/pritunl-monitoring/errortypes"
	"github.com/pritunl/pritunl-monitoring/hosts"
	"net/http"
	"os"
	"time"
)

type series struct {
	Series []*metric `json:"series"`
}

type metric struct {
	Metric string          `json:"metric"`
	Points [][]interface{} `json:"points"`
	Type   string          `json:"type"`
	Host   string          `json:"host"`
}

func Update() (err error) {
	db := database.GetDatabase()
	defer db.Close()

	data := &series{
		Series: []*metric{},
	}

	host, err := hosts.GetHost(db, os.Getenv("HOST_ID"))
	if err != nil {
		return
	}

	if host.Status != "online" {
		return
	}

	curTime := time.Now().Unix()

	fmt.Println(curTime)

	mtric := &metric{
		Metric: "pritunl.cpu_usage",
		Points: [][]interface{}{
			[]interface{}{
				curTime, host.CpuUsage * 100,
			},
		},
		Type: "gauge",
		Host: host.Name,
	}
	data.Series = append(data.Series, mtric)

	mtric = &metric{
		Metric: "pritunl.mem_usage",
		Points: [][]interface{}{
			[]interface{}{
				curTime, host.MemUsage * 100,
			},
		},
		Type: "gauge",
		Host: host.Name,
	}
	data.Series = append(data.Series, mtric)

	mtric = &metric{
		Metric: "pritunl.running_servers",
		Points: [][]interface{}{
			[]interface{}{
				curTime, host.ServerCount,
			},
		},
		Type: "gauge",
		Host: host.Name,
	}
	data.Series = append(data.Series, mtric)

	mtric = &metric{
		Metric: "pritunl.connected_devices",
		Points: [][]interface{}{
			[]interface{}{
				curTime, host.DeviceCount,
			},
		},
		Type: "gauge",
		Host: host.Name,
	}
	data.Series = append(data.Series, mtric)

	mtric = &metric{
		Metric: "pritunl.thread_count",
		Points: [][]interface{}{
			[]interface{}{
				curTime, host.ThreadCount,
			},
		},
		Type: "gauge",
		Host: host.Name,
	}
	data.Series = append(data.Series, mtric)

	jsonData, err := json.Marshal(data)
	if err != nil {
		err = errortypes.ParseError{
			errors.Wrap(err, "datadog: Failed to parse series data"),
		}
		return
	}

	resp, err := http.Post(
		fmt.Sprintf(
			"https://app.datadoghq.com/api/v1/series?api_key=%s",
			os.Getenv("DATADOG_API_KEY")),
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		err = errortypes.RequestError{
			errors.New("datadog: Failed to post series"),
		}
		return
	}

	return
}

func Start() (err error) {
	Update()

	for {
		time.Sleep(10 * time.Second)
		err := Update()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("datadog: Update error")
		}
	}

	return
}
