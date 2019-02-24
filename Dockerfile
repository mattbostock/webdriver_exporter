FROM golang as builder

WORKDIR /go/src/app
COPY main.go .
COPY probe.go .
COPY vendor/ vendor/

RUN ls -altr
RUN go get ./...
RUN go build

FROM robcherry/docker-chromedriver
##FROM selenium/standalone-chrome

COPY --from=builder /go/src/app/app /usr/local/bin/webdriver_exporter
COPY etc/supervisor/conf.d/exporter.conf /etc/supervisor/conf.d/exporter.conf
RUN useradd -ms /bin/bash dev

EXPOSE 9156
