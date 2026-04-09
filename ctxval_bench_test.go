package ctxval_test

import (
	"context"
	"testing"

	ctxval "github.com/selfshop-dev/lib-ctxval"
)

// BenchmarkWith measures the cost of storing a value into a fresh context layer.
func BenchmarkWith(b *testing.B) {
	ctx := context.Background()
	b.ReportAllocs()

	for b.Loop() {
		_ = ctxval.With(ctx, tString("bench"))
	}
}

// BenchmarkGet_hit measures retrieval when the value is present at the top of
// the context chain — the best-case traversal.
func BenchmarkGet_hit(b *testing.B) {
	ctx := ctxval.With(context.Background(), tString("bench"))
	b.ReportAllocs()

	for b.Loop() {
		_, _ = ctxval.Get[tString](ctx)
	}
}

// BenchmarkGet_miss measures retrieval when the type was never stored —
// the full chain must be traversed before returning false.
func BenchmarkGet_miss(b *testing.B) {
	ctx := context.Background()
	b.ReportAllocs()

	for b.Loop() {
		_, _ = ctxval.Get[tString](ctx)
	}
}

// BenchmarkOr_hit measures [Or] when the value is present.
func BenchmarkOr_hit(b *testing.B) {
	ctx := ctxval.With(context.Background(), tString("bench"))
	b.ReportAllocs()

	for b.Loop() {
		_ = ctxval.Or(ctx, tString("fallback"))
	}
}

// BenchmarkOr_miss measures [Or] when nothing is stored and the fallback
// is returned.
func BenchmarkOr_miss(b *testing.B) {
	ctx := context.Background()
	b.ReportAllocs()

	for b.Loop() {
		_ = ctxval.Or(ctx, tString("fallback"))
	}
}

// BenchmarkGet_deepChain measures lookup cost when the target value is buried
// under several unrelated context layers — the worst-case traversal path.
func BenchmarkGet_deepChain(b *testing.B) {
	ctx := ctxval.With(context.Background(), tString("deep"))
	for range 10 {
		//nolint:fatcontext // noise layers of a different type
		ctx = ctxval.With(ctx, tInt(0))
		ctx = context.WithValue(ctx, struct{ k int }{}, 0) // opaque key for benchmarking depth
	}
	b.ReportAllocs()

	for b.Loop() {
		_, _ = ctxval.Get[tString](ctx)
	}
}
