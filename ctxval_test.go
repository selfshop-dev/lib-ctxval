package ctxval_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ctxval "github.com/selfshop-dev/lib-ctxval"
)

type (
	tString  string
	tInt     int
	tFloat64 float64
	tStruct  struct{ V int }
)

func TestWith_Get_roundtrip(t *testing.T) {
	t.Parallel()

	ctx := ctxval.With(context.Background(), tString("hello"))

	got, ok := ctxval.Get[tString](ctx)

	require.True(t, ok)
	assert.Equal(t, tString("hello"), got)
}

func TestGet_missing_returnsFalse(t *testing.T) {
	t.Parallel()

	got, ok := ctxval.Get[tString](context.Background())

	assert.False(t, ok)
	assert.Equal(t, tString(""), got)
}

func TestWith_Get_zeroValue(t *testing.T) {
	t.Parallel()

	// Explicitly storing the zero value must be distinguishable from absence.
	ctx := ctxval.With(context.Background(), tString(""))

	got, ok := ctxval.Get[tString](ctx)

	require.True(t, ok, "zero value was stored but Get reported missing")
	assert.Equal(t, tString(""), got)
}

func TestWith_Get_independentTypes(t *testing.T) {
	t.Parallel()

	ctx := ctxval.With(context.Background(), tString("str"))
	ctx = ctxval.With(ctx, tInt(42))
	ctx = ctxval.With(ctx, tFloat64(3.14))

	s, okS := ctxval.Get[tString](ctx)
	i, okI := ctxval.Get[tInt](ctx)
	f, okF := ctxval.Get[tFloat64](ctx)

	require.True(t, okS)
	require.True(t, okI)
	require.True(t, okF)
	assert.Equal(t, tString("str"), s)
	assert.Equal(t, tInt(42), i)

	//nolint:testifylint // exact comparison is safe in controlled test values
	assert.Equal(t, tFloat64(3.14), f)
}

func TestWith_shadow_latestWins(t *testing.T) {
	t.Parallel()

	ctx := ctxval.With(context.Background(), tString("first"))
	ctx = ctxval.With(ctx, tString("second"))

	got, ok := ctxval.Get[tString](ctx)

	require.True(t, ok)
	assert.Equal(t, tString("second"), got)
}

func TestWith_shadow_doesNotMutateParent(t *testing.T) {
	t.Parallel()

	parent := ctxval.With(context.Background(), tString("original"))
	child := ctxval.With(parent, tString("overridden"))

	parentVal, _ := ctxval.Get[tString](parent)
	childVal, _ := ctxval.Get[tString](child)

	assert.Equal(t, tString("original"), parentVal, "parent context must not be mutated")
	assert.Equal(t, tString("overridden"), childVal)
}

func TestWith_Get_struct(t *testing.T) {
	t.Parallel()

	want := tStruct{V: 7}
	ctx := ctxval.With(context.Background(), want)

	got, ok := ctxval.Get[tStruct](ctx)

	require.True(t, ok)
	assert.Equal(t, want, got)
}

func TestWith_Get_pointer(t *testing.T) {
	t.Parallel()

	want := &tStruct{V: 99}
	ctx := ctxval.With(context.Background(), want)

	got, ok := ctxval.Get[*tStruct](ctx)

	require.True(t, ok)
	assert.Same(t, want, got)
}

func TestWith_Get_nilPointer(t *testing.T) {
	t.Parallel()

	// A typed nil (*tStruct) is a non-nil interface value — the type assertion
	// inside Get succeeds, so ok must be true.
	var p *tStruct
	ctx := ctxval.With(context.Background(), p)

	got, ok := ctxval.Get[*tStruct](ctx)

	require.True(t, ok, "typed nil pointer should be retrievable")
	assert.Nil(t, got)
}

func TestMust_present_returnsValue(t *testing.T) {
	t.Parallel()

	ctx := ctxval.With(context.Background(), tInt(5))

	assert.Equal(t, tInt(5), ctxval.Must[tInt](ctx))
}

func TestMust_missing_returnsZero(t *testing.T) {
	t.Parallel()

	assert.Equal(t, tInt(0), ctxval.Must[tInt](context.Background()))
}

func TestMust_missing_doesNotPanic(t *testing.T) {
	t.Parallel()

	assert.NotPanics(t, func() {
		_ = ctxval.Must[tString](context.Background())
	})
}

func TestOr_present_returnsStored(t *testing.T) {
	t.Parallel()

	ctx := ctxval.With(context.Background(), tString("stored"))

	assert.Equal(t, tString("stored"), ctxval.Or(ctx, tString("fallback")))
}

func TestOr_missing_returnsFallback(t *testing.T) {
	t.Parallel()

	assert.Equal(t, tString("fallback"), ctxval.Or(context.Background(), tString("fallback")))
}

func TestOr_storedZeroValue_returnsZero_notFallback(t *testing.T) {
	t.Parallel()

	// The stored zero value is legitimate — Or must return it, not the fallback.
	ctx := ctxval.With(context.Background(), tString(""))

	assert.Equal(t, tString(""), ctxval.Or(ctx, tString("fallback")), "stored zero value should win over fallback")
}

func TestConcurrentAccess(t *testing.T) {
	t.Parallel()

	ctx := ctxval.With(context.Background(), tString("shared"))

	const goroutines = 64

	ready := make(chan struct{})
	done := make(chan struct{}, goroutines*2)

	for range goroutines {
		go func() {
			<-ready
			_, _ = ctxval.Get[tString](ctx)
			done <- struct{}{}
		}()
		go func() {
			<-ready
			_ = ctxval.With(ctx, tString("writer"))
			done <- struct{}{}
		}()
	}

	close(ready)
	for range goroutines * 2 {
		<-done
	}
}
