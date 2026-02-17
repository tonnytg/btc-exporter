# BTC exporter

Simple BTC exporter to use with Prometheus and Grafana


# Prometheus

  - job_name: btc-exporter
    static_configs:
      - targets:
        - btc-exporter.finance.svc.cluster.local:8081


# Grafana

Prometheus Query:

    btc_price{currency="usd"}
