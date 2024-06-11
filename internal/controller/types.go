package controller

import "math/big"

type RateRequest struct {
	TokenIn  string
	TokenOut string
}

type RateResponse struct {
	Rate *big.Float
}
