
.PHONY: build test run

build:
	docker build -t nouchka/webdriver_exporter .

run:
	docker run --rm --shm-size=512M \
		-p 9156:9156 \
		--name webdriver_exporter \
		nouchka/webdriver_exporter

test: build run

test-metrics:
	docker exec webdriver_exporter wget -O- http://localhost:9156/probe?target=https://prometheus.io/
