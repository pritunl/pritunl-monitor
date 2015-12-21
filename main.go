package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-monitoring/handlers"
	"github.com/pritunl/pritunl-monitoring/metrics"
	"net/http"
	"os"
	"time"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	handlers.Register(router)

	server := &http.Server{
		Addr:           ":" + os.Getenv("PROMETHEUS_PORT"),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 4096,
	}

	e := metrics.Update()
	if e != nil {
		panic(e)
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
