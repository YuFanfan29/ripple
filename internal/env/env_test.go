package env_test

import (
	"testing"

	"github.com/YuFanfan29/ripple/internal/env"
)

func TestStringReturnsValueWhenSet(t *testing.T) {
	t.Setenv("RIPPLE_TEST_VALUE", "from-env")

	if got := env.String("RIPPLE_TEST_VALUE", "fallback"); got != "from-env" {
		t.Errorf("String() = %q, want %q", got, "from-env")
	}
}

func TestStringReturnsFallbackWhenEmpty(t *testing.T) {
	t.Setenv("RIPPLE_TEST_VALUE", "")

	if got := env.String("RIPPLE_TEST_VALUE", "fallback"); got != "fallback" {
		t.Errorf("String() = %q, want %q", got, "fallback")
	}
}
