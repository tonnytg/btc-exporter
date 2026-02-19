// cmd/crypto-exporter/main.go
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/kelseyhightower/envconfig"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"

    "myproject/internal/collector"
    "myproject/internal/client"
)

type Config struct {
    ListenAddr   string        `env:"LISTEN_ADDR" default:":8081"`
    FetchPeriod  time.Duration `env:"FETCH_PERIOD" default:"30s"`
    CoinGeckoURL string        `env:"COINGECKO_URL" default:"https://api.coingecko.com/api/v3/simple/price"`
    Coins        []string      `env:"COINS" default:"bitcoin,ethereum,solana"`
    Currencies   []string      `env:"CURRENCIES" default:"usd,brl"`
}

func main() {
    var cfg Config
    if err := envconfig.Process("", &cfg); err != nil {
        log.Fatalf("config error: %v", err)
    }

    // -------- Register metrics ----------
    prometheus.MustRegister(
        collector.PriceGauge,
        collector.ExporterUp,
        collector.ScrapeDuration,
    )
    collector.InitMetrics(cfg.Coins, cfg.Currencies)

    // -------- HTTP client ----------
    httpClient := &http.Client{Timeout: 5 * time.Second}
    api := client.NewCoingecko(httpClient, cfg.CoinGeckoURL, cfg.Coins, cfg.Currencies)

    // -------- Background fetch ----------
    ctx, stop := signal.NotifyContext(context.Background(),
        syscall.SIGINT, syscall.SIGTERM)
    defer stop()

    go func() {
        ticker := time.NewTicker(cfg.FetchPeriod)
        defer ticker.Stop()
        for {
            select {
            case <-ctx.Done():
                return
            case <-ticker.C:
                fetch(ctx, api)
            }
        }
    }()

    // -------- HTTP server ----------
    http.Handle("/metrics", promhttp.Handler())
    srv := &http.Server{Addr: cfg.ListenAddr, Handler: nil}

    go func() {
        <-ctx.Done()
        shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        _ = srv.Shutdown(shutdownCtx)
    }()

    log.Printf("crypto exporter listening on %s", cfg.ListenAddr)
    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("server error: %v", err)
    }
}

// fetch wraps the whole cycle (metrics, logging, health flag)
func fetch(ctx context.Context, api *client.Coingecko) {
    start := time.Now()
    defer func() { collector.ScrapeDuration.Observe(time.Since(start).Seconds()) }()

    prices, err := api.FetchPrices(ctx)
    if err != nil {
        collector.ExporterUp.Set(0)
        log.Printf("fetch error: %v", err)
        return
    }
    collector.ExporterUp.Set(1)

    for _, p := range prices {
        collector.PriceGauge.WithLabelValues(p.Coin, p.Currency).Set(p.Price)
    }
}

