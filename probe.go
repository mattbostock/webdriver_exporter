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
	"fmt"
	"net/http"

	"github.com/prometheus/common/log"
	"github.com/sclevine/agouti"
)

const msInSecond = 1000

// from https://www.w3.org/TR/navigation-timing/#sec-navigation-timing-interface
type navigationTimings struct {
	ConnectEnd                 float64 `json:"connectEnd"`
	ConnectStart               float64 `json:"connectStart"`
	DomComplete                float64 `json:"domComplete"`
	DomContentLoadedEventEnd   float64 `json:"domContentLoadedEventEnd"`
	DomContentLoadedEventStart float64 `json:"domContentLoadedEventStart"`
	DomInteractive             float64 `json:"domInteractive"`
	DomLoading                 float64 `json:"domLoading"`
	DomainLookupEnd            float64 `json:"domainLookupEnd"`
	DomainLookupStart          float64 `json:"domainLookupStart"`
	FetchStart                 float64 `json:"fetchStart"`
	LoadEventEnd               float64 `json:"loadEventEnd"`
	LoadEventStart             float64 `json:"loadEventStart"`
	NavigationStart            float64 `json:"navigationStart"`
	RedirectEnd                float64 `json:"redirectEnd"`
	RedirectStart              float64 `json:"redirectStart"`
	RequestStart               float64 `json:"requestStart"`
	ResponseEnd                float64 `json:"responseEnd"`
	ResponseStart              float64 `json:"responseStart"`
	SecureConnectionStart      float64 `json:"secureConnectionStart"`
	UnloadEventEnd             float64 `json:"unloadEventEnd"`
	UnloadEventStart           float64 `json:"unloadEventStart"`
}

func probe(target string, w http.ResponseWriter) bool {
	// new page to ensure a clean session (no caching)
	options := []agouti.Option{
		agouti.Desired(agouti.NewCapabilities().Browser("chrome").With("javascriptEnabled")),
		agouti.Timeout(5),
	}
	page, err := driver.NewPage(options...)
	if err != nil {
		log.Error(err)
		return false
	}
	defer page.Destroy()

	err = page.Navigate(target)
	if err != nil {
		log.Error(err)
		return false
	}
	url, err := page.URL()
	if err != nil {
		log.Error(err)
		return false
	}
	if url != target {
		log.Errorf("got unexpected URL, ensure target accounts for any redirects; got %q, expected %q", url, target)
		return false
	}

	var timings *navigationTimings
	err = page.RunScript("return window.performance.timing;", nil, &timings)
	if err != nil {
		log.Error(err)
		return false
	}

	fmt.Fprintf(w, "navigation_timing_connect_end_seconds %f\n", timings.ConnectEnd/msInSecond)
	fmt.Fprintf(w, "navigation_timing_connect_start_seconds %f\n", timings.ConnectStart/msInSecond)
	fmt.Fprintf(w, "navigation_timing_dom_complete_seconds %f\n", timings.DomComplete/msInSecond)
	fmt.Fprintf(w, "navigation_timing_dom_content_loaded_event_end_seconds %f\n", timings.DomContentLoadedEventEnd/msInSecond)
	fmt.Fprintf(w, "navigation_timing_dom_content_loaded_event_start_seconds %f\n", timings.DomContentLoadedEventStart/msInSecond)
	fmt.Fprintf(w, "navigation_timing_dom_interactive_seconds %f\n", timings.DomInteractive/msInSecond)
	fmt.Fprintf(w, "navigation_timing_dom_loading_seconds %f\n", timings.DomLoading/msInSecond)
	fmt.Fprintf(w, "navigation_timing_domain_lookup_end_seconds %f\n", timings.DomainLookupEnd/msInSecond)
	fmt.Fprintf(w, "navigation_timing_domain_lookup_start_seconds %f\n", timings.DomainLookupStart/msInSecond)
	fmt.Fprintf(w, "navigation_timing_fetch_start_seconds %f\n", timings.FetchStart/msInSecond)
	fmt.Fprintf(w, "navigation_timing_load_event_end_seconds %f\n", timings.LoadEventEnd/msInSecond)
	fmt.Fprintf(w, "navigation_timing_load_event_start_seconds %f\n", timings.LoadEventStart/msInSecond)
	fmt.Fprintf(w, "navigation_timing_redirect_end_seconds %f\n", timings.RedirectEnd/msInSecond)
	fmt.Fprintf(w, "navigation_timing_redirect_start_seconds %f\n", timings.RedirectStart/msInSecond)
	fmt.Fprintf(w, "navigation_timing_request_start_seconds %f\n", timings.RequestStart/msInSecond)
	fmt.Fprintf(w, "navigation_timing_response_end_seconds %f\n", timings.ResponseEnd/msInSecond)
	fmt.Fprintf(w, "navigation_timing_response_start_seconds %f\n", timings.ResponseStart/msInSecond)
	fmt.Fprintf(w, "navigation_timing_secure_connection_start_seconds %f\n", timings.SecureConnectionStart/msInSecond)
	fmt.Fprintf(w, "navigation_timing_start_seconds %f\n", timings.NavigationStart/msInSecond)
	fmt.Fprintf(w, "navigation_timing_unload_event_end_seconds %f\n", timings.UnloadEventEnd/msInSecond)
	fmt.Fprintf(w, "navigation_timing_unload_event_start_seconds %f\n", timings.UnloadEventStart/msInSecond)

	logs, err := page.ReadAllLogs("browser")
	if err != nil {
		log.Error(err)
		return false
	}
	var warningCount, severeCount int
	for _, log := range logs {
		if log.Level == "WARNING" {
			warningCount++
		}
		if log.Level == "SEVERE" {
			severeCount++
		}
	}
	fmt.Fprintf(w, "browser_log_warning_count %d\n", warningCount)
	fmt.Fprintf(w, "browser_log_severe_count %d\n", severeCount)

	return true
}
