package chaos

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestApply_Deterministic(t *testing.T) {
	ctx := context.Background()

	t.Run("no chaos", func(t *testing.T) {
		cfg := Config{DelayPercent: 0, DelayMS: 0, ErrorPercent: 0}
		if err := cfg.Apply(ctx); err != nil {
			t.Fatalf("Apply = %v, want nil", err)
		}
	})

	t.Run("always error", func(t *testing.T) {
		cfg := Config{ErrorPercent: 100}
		if err := cfg.Apply(ctx); !errors.Is(err, ErrChaos) {
			t.Fatalf("Apply = %v, want ErrChaos", err)
		}
	})
}

func TestApply_Delay(t *testing.T) {
	ctx := context.Background()
	cfg := Config{DelayPercent: 100, DelayMS: 50}

	start := time.Now()
	if err := cfg.Apply(ctx); err != nil {
		t.Fatalf("Apply = %v, want nil", err)
	}
	if elapsed := time.Since(start); elapsed < 50*time.Millisecond {
		t.Fatalf("elapsed = %v, want >= 50ms", elapsed)
	}
}

func TestApply_ContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	cfg := Config{DelayPercent: 100, DelayMS: 10_000}
	if err := cfg.Apply(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("Apply = %v, want context.Canceled", err)
	}
}

func TestFromEnv(t *testing.T) {
	t.Setenv("PAYMENT_DELAY_PERCENT", "50")
	t.Setenv("PAYMENT_DELAY_MS", "100")
	t.Setenv("PAYMENT_ERROR_PERCENT", "25")

	cfg := FromEnv()
	if cfg.DelayPercent != 50 || cfg.DelayMS != 100 || cfg.ErrorPercent != 25 {
		t.Fatalf("cfg = %+v, want {50 100 25}", cfg)
	}

	t.Setenv("PAYMENT_DELAY_PERCENT", "not-a-number")
	cfg = FromEnv()
	if cfg.DelayPercent != 20 {
		t.Fatalf("invalid env should fall back to default 20, got %d", cfg.DelayPercent)
	}
}
