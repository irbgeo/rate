package rater

import (
	"context"
	"time"
)

type StartOpts struct {
	RateUpdateInterval time.Duration
	PairUpdateInterval time.Duration
}

type worker struct {
	cancel context.CancelFunc
	rate   string
}
