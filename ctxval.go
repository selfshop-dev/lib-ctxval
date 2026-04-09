package ctxval

import "context"

// key is the unexported generic context key type.
// Each distinct type T produces a unique key[T]{}, eliminating collisions
// between packages without any manual key declaration.
type key[T any] struct{}

// With returns a copy of ctx with v stored under the type T.
// Subsequent calls with the same T shadow — but do not remove — earlier
// values; the full history remains traversable via the parent context chain,
// identical to [context.WithValue] semantics.
func With[T any](ctx context.Context, v T) context.Context {
	return context.WithValue(ctx, key[T]{}, v)
}

// Get retrieves the value of type T from ctx.
// It returns the stored value and true if found, or the zero value of T
// and false if T was never stored in the chain.
func Get[T any](ctx context.Context) (T, bool) {
	v, ok := ctx.Value(key[T]{}).(T)
	return v, ok
}

// Must retrieves the value of type T from ctx.
// It returns the stored value if found, or the zero value of T if not.
// Must never panics; use [Get] when absence must be distinguished from
// a legitimately stored zero value.
func Must[T any](ctx context.Context) T {
	if v, ok := Get[T](ctx); ok {
		return v
	}
	var zero T
	return zero
}

// Or retrieves the value of type T from ctx.
// It returns the stored value if found, or fallback otherwise.
func Or[T any](ctx context.Context, fallback T) T {
	if v, ok := Get[T](ctx); ok {
		return v
	}
	return fallback
}
