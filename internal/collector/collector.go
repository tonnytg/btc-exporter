// internal/collector/collector.go
package collector

import (
    "github.com/prometheus/client_golang/prometheus"
)

var (
    PriceGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
        Name: "crypto_price",
        Help: "Crypto price by coin and currency",
    }, []string{"coin", "currency"})

    ExporterUp = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "crypto_exporter_up",
        Help: "Exporter health status (1 = up, 0 = down)",
    })

    ScrapeDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
        Name:    "crypto_exporter_scrape_duration_seconds",
        Help:    "Time spent scraping the API",
        Buckets: prometheus.ExponentialBuckets(0.1, 2, 8),
    })
)

// InitMetrics pre‑registers all label combos (avoids creating metrics at runtime)
func InitMetrics(coins, cur []string) {
    for _, coin := range coins {
        for _, c := range cur {
            PriceGauge.WithLabelValues(coin, c)
        }
    }
}


