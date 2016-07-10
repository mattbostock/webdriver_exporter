# Navigation Timing Exporter

Simulates a browser session using the [WebDriver protocol][] and exposes
[Navigation Timings][] as a [Prometheus][] exporter.

[Prometheus]: https://prometheus.io/
[Navigation Timings]: https://www.w3.org/TR/navigation-timing/
[WebDriver protocol]: https://www.w3.org/TR/webdriver/

## Alpha quality

This is alpha-quality code, without tests. Run it in Production at your own
risk.

## Example output

```
navigation_timing_secure_connection_start_seconds 1468151331.286000
navigation_timing_dom_complete_seconds 1468151332.641000
navigation_timing_dom_content_loaded_event_end_seconds 1468151332.167000
navigation_timing_load_event_end_seconds 1468151332.643000
navigation_timing_response_end_seconds 1468151331.645000
navigation_timing_redirect_start_seconds 0.000000
navigation_timing_request_start_seconds 1468151331.465000
navigation_timing_response_start_seconds 1468151331.643000
navigation_timing_dom_content_loaded_event_start_seconds 1468151332.158000
navigation_timing_dom_interactive_seconds 1468151332.158000
navigation_timing_dom_loading_seconds 1468151331.647000
navigation_timing_fetch_start_seconds 1468151331.221000
navigation_timing_connect_end_seconds 1468151331.465000
navigation_timing_connect_start_seconds 1468151331.277000
navigation_timing_domain_lookup_end_seconds 1468151331.277000
navigation_timing_unload_event_start_seconds 0.000000
navigation_timing_unload_event_end_seconds 0.000000
navigation_timing_domain_lookup_start_seconds 1468151331.271000
navigation_timing_load_event_start_seconds 1468151332.641000
navigation_timing_navigation_start_seconds 1468151330.960000
navigation_timing_redirect_end_seconds 0.000000
probe_duration_seconds 2.596026
probe_success 1
```

## Building and running

### Prerequisites

You'll need [chromedriver][]:

    # On Mac OS X using Homebrew
    brew install chromedriver

To run Navigation Timing Exporter on a server with a headless Chrome browser,
you'll need something like [xvfb][].

[chromedriver]: https://sites.google.com/a/chromium.org/chromedriver/
[xvfb]: https://www.x.org/archive/X11R7.6/doc/man/man1/Xvfb.1.xhtml

### Building locally

    go build
    ./navigation_timing_exporter <flags>

Visiting [http://localhost:9116/probe?target=https://prometheus.io/](http://localhost:9116/probe?target=https://prometheus.io/)
will return metrics from the Navigation Timing API for prometheus.io.

## Prometheus Configuration

The Navigation Timing Exporter needs to be passed the target as a parameter,
this can be done with relabelling.

Example configuration:

```yaml
scrape_configs:
  - job_name: 'navigation_timing'
    metrics_path: /probe
    target_groups:
      - targets:
        - https://prometheus.io/   # Target to probe
    relabel_configs:
      - source_labels: [__address__]
        regex: (.*)(:80)?
        target_label: __param_target
        replacement: ${1}
      - source_labels: [__param_target]
        regex: (.*)
        target_label: instance
        replacement: ${1}
      - source_labels: []
        regex: .*
        target_label: __address__
        replacement: 127.0.0.1:9116  # Navigation Timing Exporter
```

## Limitations

### Chromedriver only (currently)

Only Chromedriver is supported currently. Adding support for other webdrivers
(e.g. Selenium, phantomjs) should be relatively trivial. Pull requests are
encouraged.

At the time of writing, PhantomJS did not support probing of sites with strict
[Content Security Policies][]; see ariya/phantomjs#13114.

[Content Security Policy]: https://www.w3.org/TR/CSP1/

### Requires exact URL as target

To ensure that the timings returned are for the target requested, the page URL
is strictly matched against the requested target.

For example, if you try to probe `https://example.com/`, the probe will fail if
the page redirects to `https://example.com/foo`. Similarly, make sure that you
use `https` as the URL scheme if your target enforces HTTPS.

This requirement could be loosened if there's enough demand for it.

### Probes are slow to initialise

We start a new [Chromedriver][] session for each probe to ensure that the
cache, cookies and local storage are clean. If we could retain the session and
clear the cache, the time to complete a probe would be significantly reduced.

Currently each probe takes around 2 seconds.
