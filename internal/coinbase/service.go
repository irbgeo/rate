package coinbase

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const coinbaseAPIURL = "https://api.coinbase.com/v2/prices/%s-%s/spot"

type service struct {
}

func NewService() *service {
	return &service{}
}

func (s *service) Rate(ctx context.Context, tokenIn, tokenOut string) (string, error) {
	url := fmt.Sprintf(coinbaseAPIURL, tokenIn, tokenOut)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get response from Coinbase API: %s", resp.Status)
	}

	var coinbaseResponse CoinbaseResponse
	err = json.NewDecoder(resp.Body).Decode(&coinbaseResponse)
	if err != nil {
		return "", err
	}

	return coinbaseResponse.Data.Amount, nil
}
