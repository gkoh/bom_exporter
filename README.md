# bom_exporter - Prometheus exporter for the Australian Bureau of Meteorology

A Prometheus compatible exporter for forecast and observation data published by
the Australian Bureau of Meteorology.

The exporter adapts the forecast and observation data feeds (from the anonymous
FTP BOM server) into promscrape compatible format.

## Build
```
go build cmd/bom_exporter.go
```

A binary called `bom_exporter` should be compiled.

## Test
```
go test ./...
```

## Scraping Data
`bom_exporter` will serve forecast scrape data on the following URL:
```
http://<server>:8080/metrics?<product_id>
```
Where the product identifier can be obtained from:
http://www.bom.gov.au/catalogue/anon-ftp.shtml

The following product identifier types are currently supported:
- forecast
- observations

## Scrape Configuration
Data is retrieved from the BoM on request, thus the configured scrape interval
controls the frequency at which the BoM FTP server is accessed.
It is recommended to match the scrape interval with the product being accessed:
- observations - scrape every 5 minutes, updated every 10 minutes
- forecast - scrape every 1-6 hours, updated every 12(?) hours

### Example Scrape Single Product
Following is the configuration snippet to scrape the Sydney city forecast
(IDN10064) hourly.
```
  - job_name: bom_forecast
    scrape_interval: 1h
    metrics_path: /metrics
    params:
      id: ['IDN10064']
    static_configs:
      - targets: ['localhost:8080']
```

### Example Scrape Multiple Products
Following is the configuration snippet to scrape all the state observations
every 5 minutes:
```
  - job_name: bom_observations
    scrape_interval: 5m
    metrics_path: /metrics
    static_configs:
      - targets:
        - 'IDS60920'
        - 'IDN60920'
        - 'IDD60920'
        - 'IDQ60920'
        - 'IDT60920'
        - 'IDV60920'
        - 'IDW60920'
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_id
      - source_labels: [__param_id]
        target_label: instance
      - target_label: __address__
        replacement: localhost:8080
```

## Motivation
I've always wanted to have longer term climate data available with a user
interface that I have familiarity (Grafana).

## Known Issues/Future Ideas
- Using a scrape interval >5 minutes results in Prometheus staleness
  - Cache the data and timestamps internally
  - Disconnect the external scrape interval from the data retrieval and use the
     'next issue time' to intelligently schedule the next FTP retrieval.
- Support other products, eg. tides
- Release a Docker image
  - Look into goreleaser
- Improve test coverage
- Add continuous integration
