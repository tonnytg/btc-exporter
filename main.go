package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type CoinPrice struct {
	USD float64 `json:"usd"`
	BRL float64 `json:"brl"`
}

type Response struct {
	Bitcoin  CoinPrice `json:"bitcoin"`
	Ethereum CoinPrice `json:"ethereum"`
	Solana   CoinPrice `json:"solana"`
}

var (
	cryptoPrice = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "crypto_price",
			Help: "Crypto price by coin and currency",
		},
		[]string{"coin", "currency"},
	)

	httpClient = &http.Client{
		Timeout: 5 * time.Second,
	}
)

func fetchPrices() {
	url := "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin,ethereum,solana&vs_currencies=usd,brl"

	resp, err := httpClient.Get(url)
	if err != nil {
		log.Println("error fetching prices:", err)
		return
	}
	defer resp.Body.Close()

	var data Response
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Println("error decoding JSON:", err)
		return
	}

	// BTC
	cryptoPrice.WithLabelValues("bitcoin", "usd").Set(data.Bitcoin.USD)
	cryptoPrice.WithLabelValues("bitcoin", "brl").Set(data.Bitcoin.BRL)

	// ETH
	cryptoPrice.WithLabelValues("ethereum", "usd").Set(data.Ethereum.USD)
	cryptoPrice.WithLabelValues("ethereum", "brl").Set(data.Ethereum.BRL)

	// SOL
	cryptoPrice.WithLabelValues("solana", "usd").Set(data.Solana.USD)
	cryptoPrice.WithLabelValues("solana", "brl").Set(data.Solana.BRL)
}

func main() {
	prometheus.MustRegister(cryptoPrice)

	go func() {
		for {
			fetchPrices()
			time.Sleep(30 * time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())

	log.Println("Crypto exporter running on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
