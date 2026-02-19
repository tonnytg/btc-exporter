// internal/client/coingecko.go
package client

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    "strings"

    "myproject/internal/model"
)

type Coingecko struct {
    http   *http.Client
    base   string
    coins  []string
    cur    []string
}

func NewCoingecko(httpClient *http.Client, base string, coins, cur []string) *Coingecko {
    return &Coingecko{http: httpClient, base: base, coins: coins, cur: cur}
}

func (c *Coingecko) buildURL() string {
    v := url.Values{}
    v.Set("ids", strings.Join(c.coins, ","))
    v.Set("vs_currencies", strings.Join(c.cur, ","))
    return fmt.Sprintf("%s?%s", c.base, v.Encode())
}

func (c *Coingecko) FetchPrices(ctx context.Context) ([]*model.PriceMetric, error) {
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.buildURL(), nil)
    if err != nil {
        return nil, err
    }

    resp, err := c.http.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
    }

    var raw map[string]map[string]float64
    if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
        return nil, err
    }

    // flatten para o collector
    var out []*model.PriceMetric
    for coin, curMap := range raw {
        for cur, price := range curMap {
            out = append(out, &model.PriceMetric{
                Coin:     coin,
                Currency: cur,
                Price:    price,
            })
        }
    }
    return out, nil
}

