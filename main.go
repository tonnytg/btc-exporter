package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Response struct {
	Bitcoin struct {
		USD float64 `json:"usd"`
		BRL float64 `json:"brl"`
	} `json:"bitcoin"`
}

var (
	btcPrice = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "btc_price",
			Help: "Bitcoin price by currency",
		},
		[]string{"currency"},
	)

	httpClient = &http.Client{
		Timeout: 5 * time.Second,
	}
)

func fetchBTC() {
	resp, err := httpClient.Get("https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd,brl")
	if err != nil {
		log.Println("error fetching BTC:", err)
		return
	}
	defer resp.Body.Close()

	var data Response
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Println("error decoding JSON:", err)
		return
	}

	btcPrice.WithLabelValues("usd").Set(data.Bitcoin.USD)
	btcPrice.WithLabelValues("brl").Set(data.Bitcoin.BRL)
}

func main() {
	prometheus.MustRegister(btcPrice)

	go func() {
		for {
			fetchBTC()
			time.Sleep(30 * time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())

	log.Println("BTC exporter running on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

