package metrics

import (
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "net/http"
)

func InitMetrics() {
    go ExposeMetrics()
}

func ExposeMetrics() {
    http.Handle("/metrics", promhttp.Handler())
    _ = http.ListenAndServe(":2112", nil)
}

