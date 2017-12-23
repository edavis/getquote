package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Price struct {
	Coin struct {
		BTC    float64 `json:"btc"`
		Name   string  `json:"name"`
		Rank   int64   `json:"rank"`
		Ticker string  `json:"ticker"`
		USD    float64 `json:"usd"`
	} `json:"coin"`
}

func getPrice(coin string) (*Price, error) {
	var url = fmt.Sprintf("https://coinbin.org/%s", coin)
	var resp, err = http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var price = new(Price)
	var decoder = json.NewDecoder(resp.Body)
	if err := decoder.Decode(price); err != nil {
		return nil, err
	}
	return price, nil
}

func main() {
	var now = time.Now().Format("2006/01/02 15:04:05")
	// TODO(ejd): put this in an envvar or something
	var coins = []string{"BTC", "BCH", "ETH", "LTC"}

	var pricedb, err = os.OpenFile(os.Getenv("LEDGER_PRICE_DB"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err.Error())
	}
	defer pricedb.Close()

	for _, coin := range coins {
		var price, err = getPrice(coin)
		if err != nil {
			fmt.Fprintf(pricedb, "# %s: could not obtain price info for %s\n", now, coin)
			continue
		}
		fmt.Fprintf(pricedb, "P %s %s $%.2f\n", now, coin, price.Coin.USD)
	}
}
