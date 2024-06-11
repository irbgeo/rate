package rater

import (
	"context"
	"log/slog"
	"time"
)

func (s *service) updateWorkers(tokens []string) error {
	pairs := buildPairNames(tokens)

	s.worker.Range(func(key, value interface{}) bool {
		pairName := key.(string)
		if _, ok := pairs[pairName]; ok {
			delete(pairs, pairName)
			return true
		}

		u := value.(worker)
		u.cancel()

		s.worker.Delete(pairName)
		return true

	})

	for pairName := range pairs {
		tokenA, tokenB := parsePairName(pairName)

		err := s.startWorker(tokenA, tokenB)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *service) startWorker(tokenA, tokenB string) error {
	err := s.updateRate(tokenA, tokenB)
	if err != nil {
		return err
	}

	go func() {
		ticker := time.NewTicker(s.rateUpdateInterval)
		defer ticker.Stop()

		for {
			select {
			case <-s.ctx.Done():
				return
			case <-ticker.C:
				err := s.updateRate(tokenA, tokenB)
				if err != nil {
					slog.Error("failed to update rate", "tokenA", tokenA, "tokenB", tokenB, "err", err)
				}
			}
		}
	}()
	return nil
}

func (s *service) updateRate(tokenA, tokenB string) error {
	ctx, cancel := context.WithTimeout(s.ctx, requestTimeout)
	defer cancel()

	rate, err := s.rateSource.Rate(ctx, tokenA, tokenB)
	if err != nil {
		return err
	}

	pairName, _ := buildPairName(tokenA, tokenB)

	_, updaterCancel := context.WithCancel(s.ctx)

	worker := worker{
		cancel: updaterCancel,
		rate:   rate,
	}

	s.worker.Store(pairName, worker)
	return nil
}
