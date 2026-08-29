package chaos

import (
	"context"
	"errors"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

// ErrChaos is returned when a forced failure is triggered.
var ErrChaos = errors.New("chaos induced payment failure")

// Config controls artificial latency and error injection for the Payment service.
type Config struct {
	DelayPercent int
	DelayMS      int
	ErrorPercent int
}

// FromEnv reads the chaos configuration from environment variables,
// falling back to the lab defaults from PRD §24.
func FromEnv() Config {
	return Config{
		DelayPercent: envInt("PAYMENT_DELAY_PERCENT", 20),
		DelayMS:      envInt("PAYMENT_DELAY_MS", 2000),
		ErrorPercent: envInt("PAYMENT_ERROR_PERCENT", 10),
	}
}

var (
	rngMu sync.Mutex
	rng   = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func roll() int {
	rngMu.Lock()
	n := rng.Intn(100)
	rngMu.Unlock()
	return n
}

// Apply introduces latency and/or a forced error according to the config.
// It is deterministic at the extremes (0% = never, 100% = always) so tests
// can rely on it.
func (c Config) Apply(ctx context.Context) error {
	if c.ErrorPercent >= 100 {
		return ErrChaos
	}
	if c.DelayPercent <= 0 && c.ErrorPercent <= 0 {
		return nil
	}

	if c.DelayPercent > 0 && roll() < c.DelayPercent {
		select {
		case <-time.After(time.Duration(c.DelayMS) * time.Millisecond):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	if c.ErrorPercent > 0 && roll() < c.ErrorPercent {
		return ErrChaos
	}
	return nil
}

func envInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
