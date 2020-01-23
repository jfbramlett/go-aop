package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func InitMetrics(cfg MetricsConfig) {
	go ExposeMetrics(cfg)
}

func ExposeMetrics(cfg MetricsConfig) {
	http.Handle(cfg.URL, promhttp.Handler())
	_ = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil)
}
