package prometheus

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/dropbox/godropbox/errors"
	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-monitoring/errortypes"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"os"
	"time"
)

func Limiter(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1000000)
}

func Recovery(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithFields(logrus.Fields{
				"error": errors.New(fmt.Sprintf("%s", r)),
			}).Error("prometheus: Handler panic")
			c.Writer.WriteHeader(http.StatusInternalServerError)
		}
	}()

	c.Next()
}

func Register(engine *gin.Engine) {
	engine.Use(Limiter)
	engine.Use(Recovery)

	engine.GET("/metrics", metricsHandler())
}

func metricsHandler() gin.HandlerFunc {
	prmHandler := prometheus.Handler()

	return func(c *gin.Context) {
		prmHandler.ServeHTTP(c.Writer, c.Request)
	}
}

func Start() (err error) {
	Update()

	go func() {
		for {
			time.Sleep(10 * time.Second)
			err := Update()
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("prometheus: Update error")
			}
		}
	}()

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	Register(router)

	server := &http.Server{
		Addr:           ":" + os.Getenv("PROMETHEUS_PORT"),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 4096,
	}

	err = server.ListenAndServe()
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "prometheus: Server error"),
		}
		return
	}

	return
}
