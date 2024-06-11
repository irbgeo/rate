package rater

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"sync"
	"time"
)

var (
	requestTimeout = 5 * time.Second
)

type service struct {
	ctx    context.Context
	cancel context.CancelFunc
	worker sync.Map

	storage    storage
	rateSource rateSource

	rateUpdateInterval time.Duration
	pairUpdateInterval time.Duration
}

type storage interface {
	Tokens(ctx context.Context) ([]string, error)
}

type rateSource interface {
	Rate(ctx context.Context, tokenIn, tokenOut string) (string, error)
}

func NewService(s storage, src rateSource) *service {
	svc := &service{
		storage:    s,
		rateSource: src,
	}

	svc.ctx, svc.cancel = context.WithCancel(context.Background())
	return svc
}

func (s *service) Start(opts StartOpts) error {
	s.rateUpdateInterval = opts.RateUpdateInterval
	s.pairUpdateInterval = opts.PairUpdateInterval

	err := s.runtime()
	return err
}

func (s *service) Stop() error {
	s.cancel()
	return nil
}

func (s *service) Get(ctx context.Context, tokenIn, tokenOut string) (*big.Float, error) {
	reverse := false
	directPairName, reversePairName := buildPairName(tokenIn, tokenOut)

	ww, ok := s.worker.Load(directPairName)
	if !ok {
		ww, reverse = s.worker.Load(reversePairName)
		if !reverse {
			return nil, fmt.Errorf("rate for pair %s not found", directPairName)
		}
	}

	w := ww.(worker)

	rate, ok := new(big.Float).SetString(w.rate)
	if !ok {
		return nil, fmt.Errorf("rate for pair %s not found", directPairName)
	}

	if reverse {
		rate = new(big.Float).Quo(big.NewFloat(1), rate)
	}

	return rate, nil
}

func (s *service) runtime() error {
	err := s.updatePairs()
	if err != nil {
		return err
	}

	go func() {
		ticker := time.NewTicker(s.pairUpdateInterval)
		defer ticker.Stop()

		for {
			select {
			case <-s.ctx.Done():
			case <-ticker.C:
				err := s.updatePairs()
				if err != nil {
					slog.Error("failed to update pairs", slog.Any("error", err))
				}
			}
		}
	}()

	return nil
}

func (s *service) updatePairs() error {
	ctx, cancel := context.WithTimeout(s.ctx, requestTimeout)
	defer cancel()

	tokens, err := s.storage.Tokens(ctx)
	if err != nil {
		return err
	}

	if len(tokens) > 1 {
		err = s.updateWorkers(tokens)
		if err != nil {
			return err
		}
	}

	return nil
}
