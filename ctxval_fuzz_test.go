package ctxval

import (
	"context"
	"testing"
)

// FuzzGet exercises the With → Get round-trip against arbitrary string values.
//
// Invariants checked on every corpus entry:
//  1. Get always returns the exact value that was stored with With.
//  2. Get always returns ok == true for a value that was stored.
//  3. A value stored under one type is never visible under a different type.
//
// Run the fuzzer:
//
//	go test -fuzz=FuzzGet -fuzztime=60s
//
// Run the seed corpus only (CI):
//
//	go test -run=FuzzGet
func FuzzGet(f *testing.F) {
	f.Add("")
	f.Add("hello")
	f.Add("unicode: 日本語")
	f.Add("\x00\xff\xfe")

	type slotA string
	type slotB string

	f.Fuzz(func(t *testing.T, s string) {
		ctx := With(context.Background(), slotA(s))

		// Invariant 1 + 2: round-trip fidelity.
		got, ok := Get[slotA](ctx)
		if !ok {
			t.Fatalf("Get[slotA](ctx) ok = false, want true for input %q", s)
		}
		if string(got) != s {
			t.Fatalf("Get[slotA](ctx) = %q, want %q", got, s)
		}

		// Invariant 3: slotB was never stored — must not leak slotA's value.
		_, okB := Get[slotB](ctx)
		if okB {
			t.Fatalf("Get[slotB](ctx) ok = true, want false — type isolation violated for input %q", s)
		}
	})
}

// FuzzOr exercises the Or fallback contract against arbitrary string values.
//
// Invariants:
//  1. When a value is stored, Or returns that value — not the fallback.
//  2. When nothing is stored, Or returns the fallback unchanged.
func FuzzOr(f *testing.F) {
	f.Add("", "fallback")
	f.Add("stored", "fallback")
	f.Add("stored", "")
	f.Add("\x00", "\xff")

	type slot string

	f.Fuzz(func(t *testing.T, stored, fallback string) {
		// Invariant 1: stored value wins over fallback.
		ctx := With(context.Background(), slot(stored))
		if got := Or(ctx, slot(fallback)); string(got) != stored {
			t.Fatalf("Or(ctx, %q) = %q, want stored %q", fallback, got, stored)
		}

		// Invariant 2: empty context → fallback returned unchanged.
		if got := Or(context.Background(), slot(fallback)); string(got) != fallback {
			t.Fatalf("Or(emptyCtx, %q) = %q, want fallback", fallback, got)
		}
	})
}
