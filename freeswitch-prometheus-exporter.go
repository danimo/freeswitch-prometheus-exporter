package main

import (
	"log"
	"log/syslog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tomponline/fsclient/fsclient"
)

var fs *fsclient.Client

var (
	activeCalls = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freeswitch_active_calls",
		Help: "The total number of active calls",
	})
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
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

	fs_host := getEnv("FS_HOST", "127.0.0.1:8021")
	fs_password := getEnv("FS_PASSWORD", "ClueCon")
	fs = fsclient.NewClient(fs_host, fs_password, filters, subs, initFunc)

	go statsPoller()
	http.Handle("/metrics", promhttp.Handler())
	listen := getEnv("LISTEN", ":2112")
	http.ListenAndServe(listen, nil)
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
