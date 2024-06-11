package http

import "errors"

var (
	errTokenInIsRequired  = errors.New("tokenIn is required")
	errTokenOutIsRequired = errors.New("tokenOut is required")
)
