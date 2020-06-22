package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/thelastguardian/tplinkexporter/clients"
	"github.com/thelastguardian/tplinkexporter/collectors"
)

func main() {
	var (
		host     = kingpin.Flag("host", "Host of target tplink easysmart switch.").Required().String()
		username = kingpin.Flag("username", "Username for switch GUI login").Default("admin").String()
		password = kingpin.Flag("password", "Password for switch GUI login").Required().String()
	)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	tplinkSwitch := clients.NewTPLinkSwitch(*host, *username, *password)
	trafficCollector := collectors.NewTrafficCollector("tplinkexporter", tplinkSwitch)
	prometheus.MustRegister(trafficCollector)
	http.Handle("/metrics", promhttp.Handler())
	log.Info("Beginning to serve on port :9717")
	log.Fatal(http.ListenAndServe(":9717", nil))
}
