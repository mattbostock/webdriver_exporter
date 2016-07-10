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
	"unicode"

	"github.com/prometheus/common/log"
	"github.com/sclevine/agouti"
)

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

	var timings map[string]float64
	err = page.RunScript("return window.performance.timing;", nil, &timings)
	if err != nil {
		log.Error(err)
		return false
	}

	for k, v := range timings {
		fmt.Fprintf(w, "navigation_timing_%s_seconds %f\n", snakeCase(k), v/1000)
	}

	return true
}

func snakeCase(in string) string {
	runes := []rune(in)
	length := len(runes)

	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}

	return string(out)
}
