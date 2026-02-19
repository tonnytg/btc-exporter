// internal/model/model.go
package model

type CoinPrice struct {
    USD float64 `json:"usd"`
    BRL float64 `json:"brl"`
}

type Response map[string]CoinPrice // dinamiza: {"bitcoin":{"usd":...,"brl":...}, ...}

