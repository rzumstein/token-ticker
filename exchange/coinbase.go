package exchange

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// https://api.pro.coinbase.com.com/
const coinbaseBaseApi = "https://api.pro.coinbase.com/"

type coinBaseClient struct {
	exchangeBaseClient
	AccessKey string
	SecretKey string
}

type coinBaseToken struct {
	TradeID int       `json:"trade_id"`
	Price   string    `json:"price"`
	Size    string    `json:"size"`
	Bid     string    `json:"bid"`
	Ask     string    `json:"ask"`
	Volume  string    `json:"volume"`
	Time    time.Time `json:"time"`
}

type coinBaseHistoricRates [][]float64

// I don't like returning a general type here, any other better way to use the factory pattern?
func NewCoinbaseClient(httpClient *http.Client) *coinBaseClient {
	return &coinBaseClient{exchangeBaseClient: *newExchangeBase(coinbaseBaseApi, httpClient)}
}

func (client *coinBaseClient) GetName() string {
	return "Coinbase"
}

func (client *coinBaseClient) GetSymbolPrice(symbol string) (*SymbolPrice, error) {
	resp, err := client.httpGet("products/"+symbol+"/ticker", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode == 404 {
		if err := decoder.Decode(resp); err != nil {
			return nil, err
		}
		return nil, errors.New("placeholder")
	}

	var token coinBaseToken
	fmt.Println(token)
	if err := decoder.Decode(&token); err != nil {
		return nil, err
	}
	fmt.Println(token)
	price, _ := strconv.ParseFloat(token.Price, 64)

	x, _ := client.GetHistoricRates(symbol)
	fmt.Println(x)

	return &SymbolPrice{
		Symbol:           strings.Split(symbol, "-")[0],
		Price:            strconv.FormatFloat(price, 'f', 2, 64),
		Source:           client.GetName(),
		UpdateAt:         token.Time,
		PercentChange1h:  1.0, //token.PercentChange1h,
		PercentChange24h: 2.0, //token.PercentChange24h
	}, nil
}

func (client *coinBaseClient) GetHistoricRates(symbol string) (coinBaseHistoricRates, error) {
	resp, err := client.httpGet("products/"+symbol+"/candles?start=2018-10-11T12:49:00Z&end=2018-10-11T12:49:01Z", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode == 404 {
		if err := decoder.Decode(resp); err != nil {
			return nil, err
		}
		return nil, errors.New("Historic rates not found for " + symbol)
	}

	var token coinBaseHistoricRates
	if err := decoder.Decode(&token); err != nil {
		return nil, err
	}
	fmt.Println(token)

	return [][]float64{}, nil
}

func init() {
	register((&coinBaseClient{}).GetName(), func(client *http.Client) ExchangeClient {
		// Limited by type system in Go, I hate wrapper/adapter
		return NewCoinbaseClient(client)
	})
}
