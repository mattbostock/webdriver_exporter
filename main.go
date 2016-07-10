// Copyright 2016 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"github.com/sclevine/agouti"
)

const versionString = "0.0.1"

var (
	driver        = agouti.ChromeDriver()
	listenAddress = flag.String("web.listen-address", "localhost:9156", "The address to listen on for HTTP requests.")
	showVersion   = flag.Bool("version", false, "Print version information.")
)

func init() {
	version.Version = versionString
	prometheus.MustRegister(version.NewCollector("webdriver_exporter"))
}

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Fprintln(os.Stdout, version.Print("webdriver_exporter"))
		os.Exit(0)
	}

	log.Infoln("Starting webdriver_exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())

	http.Handle("/metrics", prometheus.Handler())
	http.HandleFunc("/probe",
		func(w http.ResponseWriter, r *http.Request) {
			probeHandler(w, r)
		})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
            <head><title>WebDriver Exporter</title></head>
            <body>
            <h1>WebDriver Exporter</h1>
            <p><a href="/probe?target=prometheus.io">Probe prometheus.io</a></p>
            <p><a href="/metrics">Metrics</a></p>
            </body>
            </html>`))
	})

	log.Infoln("Starting webdriver")
	err := driver.Start()
	if err != nil {
		log.Fatalf("failed to start webdriver: %s", err)
	}
	defer driver.Stop()

	log.Infoln("Listening on", *listenAddress)
	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		log.Fatalf("Error starting HTTP server: %s", err)
	}
}

func probeHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	target := params.Get("target")
	if target == "" {
		http.Error(w, "Target parameter is missing", 400)
		return
	}

	start := time.Now()
	success := probe(target, w)
	fmt.Fprintf(w, "probe_duration_seconds %f\n", float64(time.Now().Sub(start))/1e9)

	if success {
		fmt.Fprintf(w, "probe_success %d\n", 1)
	} else {
		fmt.Fprintf(w, "probe_success %d\n", 0)
	}
}
