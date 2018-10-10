package exchange

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// https://api.coinbase.com.com/v2/
const coinbaseBaseApi = "https://api.coinbase.com/v2/prices/"

type coinBaseClient struct {
	exchangeBaseClient
	AccessKey string
	SecretKey string
}

type coinBaseToken struct {
	Data struct {
		Base     string `json:"base"`
		Currency string `json:"currency"`
		Amount   string `json:"amount"`
	} `json:"data"`
}

type coinbaseCurrencyResponse struct {
	Base     string
	Currency string
	Amount   string
}

type coinbaseInvalidCurrencyResponse struct {
	id      string
	message string
}

// I don't like returning a general type here, any other better way to use the factory pattern?
func NewCoinbaseClient(httpClient *http.Client) *coinBaseClient {
	return &coinBaseClient{exchangeBaseClient: *newExchangeBase(coinbaseBaseApi, httpClient)}
}

func (client *coinBaseClient) GetName() string {
	return "Coinbase"
}

func (client *coinBaseClient) GetSymbolPrice(symbol string) (*SymbolPrice, error) {
	resp, err := client.httpGet(symbol+"/spot", nil)
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

	// bb, _ := ioutil.ReadAll(resp.Body)
	// bs := string(bb)
	// fmt.Println(bs)

	var token coinBaseToken
	fmt.Println(token)
	if err := decoder.Decode(&token); err != nil {
		return nil, err
	}
	fmt.Println(token)

	return &SymbolPrice{
		Symbol:           token.Data.Base,
		Price:            token.Data.Amount,
		Source:           client.GetName(),
		UpdateAt:         time.Unix(1, 0), //time.Unix(token.LastUpdated, 0),
		PercentChange1h:  1.0,             //token.PercentChange1h,
		PercentChange24h: 2.0,             //token.PercentChange24h
	}, nil
}

func init() {
	register((&coinBaseClient{}).GetName(), func(client *http.Client) ExchangeClient {
		// Limited by type system in Go, I hate wrapper/adapter
		return NewCoinbaseClient(client)
	})
}
