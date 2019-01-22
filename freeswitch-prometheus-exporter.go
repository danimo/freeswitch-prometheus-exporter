package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tomponline/fsclient/fsclient"
	"log"
	"log/syslog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var fs *fsclient.Client

var (
	activeCalls = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freeswitch_active_calls",
		Help: "The total number of active calls",
	})
)

func main() {
	log.SetFlags(0)
	syslogWriter, err := syslog.New(syslog.LOG_INFO, "freeswitch-prometheus-exporter")
	if err == nil {
		log.SetOutput(syslogWriter)
	}

	filters := []string{
		"Event-Name HEARTBEAT",
	}

	subs := []string{
		"HEARTBEAT",
	}

	fs = fsclient.NewClient("127.0.0.1:8021", "ClueCon", filters, subs, initFunc)

	go statsPoller()
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}

func initFunc(fsclient *fsclient.Client) {
	log.Print("Connected to freeswitch")
}

func statsPoller() {
	for {
		activeCalls.Set(getCallsCount())
		time.Sleep(2 * time.Second)
	}
}

func getCallsCount() float64 {
	callCountRes, err := fs.API("show calls count")
	if err != nil {
		log.Print("Call count API error: ", err)
		return 0
	}

	callCountRes = strings.TrimSpace(callCountRes)
	callCountParts := strings.Split(callCountRes, " ")

	callCount, err := strconv.ParseFloat(callCountParts[0], 64)
	if err != nil {
		log.Print("Call count conversion error: ", err)
	}

	return callCount
}
