# BTC exporter

Simple BTC exporter to use with Prometheus and Grafana


# Prometheus

```
  - job_name: btc-exporter
    static_configs:
      - targets:
        - btc-exporter.finance.svc.cluster.local:8081
```

# Grafana

Prometheus Query:

    crypto_price{coin="bitcoin",currency="usd"}
    crypto_price{coin="bitcoin",currency="brl"}
    crypto_price{coin="ethereum",currency="usd"}
    crypto_price{coin="ethereum",currency="brl"}
    crypto_price{coin="solana",currency="usd"}
    crypto_price{coin="solana",currency="brl"}
