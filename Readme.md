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


# How to use
    Command	            O que faz
    make ou make all    Compila (build) + roda os testes (test).
    make build	        Generate ./bin/crypto-exporter.
    make test	        Unit Tests -race e gen reports
    make fmt	        Format all codes (go fmt + goimports).
    make lint	        Exec golangci-lint.
    make vuln	        Verify vulnerabilities
    make run	        Exec binary
    make dev	        Exec “hot‑reload” (precisa de air, reflex ou entr).
    make mocks	        Regenerate files
    make docker-build	Build imagem
    make docker-push	Send imagem
    make docker-run     Run imagem
    make clean	        Remove binary
