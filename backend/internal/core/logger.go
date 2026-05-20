// Package core wires up cross-cutting infrastructure such as logging.
package core

import "go.uber.org/zap"

// NewLogger builds a zap logger appropriate for the environment.
func NewLogger(isDev bool) (*zap.Logger, error) {
	if isDev {
		return zap.NewDevelopment()
	}
	return zap.NewProduction()
}
