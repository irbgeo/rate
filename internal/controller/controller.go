package controller

import (
	"context"
	"math/big"
)

type controller struct {
	rater rater
}

type rater interface {
	Get(ctx context.Context, tokenIn, tokenOut string) (*big.Float, error)
}

func New(r rater) *controller {
	return &controller{
		rater: r,
	}
}

func (s *controller) Rate(ctx context.Context, r RateRequest) (RateResponse, error) {
	rate, err := s.rater.Get(ctx, r.TokenIn, r.TokenOut)
	if err != nil {
		return RateResponse{}, err
	}

	return RateResponse{
		Rate: rate,
	}, nil
}
