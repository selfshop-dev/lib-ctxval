package ctxval_test

import (
	"context"
	"fmt"

	ctxval "github.com/selfshop-dev/lib-ctxval"
)

// Distinct named types — one per slot in the context.
type (
	requestID string
	userID    int64
	traceID   string
)

// Example demonstrates the basic round-trip: store a value with [With],
// retrieve it with [Get].
func Example() {
	ctx := ctxval.With(context.Background(), requestID("req-42"))

	id, ok := ctxval.Get[requestID](ctx)
	fmt.Println(id, ok)

	// Output: req-42 true
}

// ExampleWith_shadow shows that storing a second value of the same type
// shadows — but does not remove — the previous one, identical to
// [context.WithValue] semantics.
func ExampleWith_shadow() {
	ctx := ctxval.With(context.Background(), requestID("first"))
	ctx = ctxval.With(ctx, requestID("second"))

	id, ok := ctxval.Get[requestID](ctx)
	fmt.Println(id, ok)

	// Output: second true
}

// ExampleWith_independentTypes shows that different types never collide:
// storing a userID does not affect the requestID slot.
func ExampleWith_independentTypes() {
	ctx := ctxval.With(context.Background(), requestID("req-1"))
	ctx = ctxval.With(ctx, userID(99))

	rid, _ := ctxval.Get[requestID](ctx)
	uid, _ := ctxval.Get[userID](ctx)
	fmt.Println(rid, uid)

	// Output: req-1 99
}

// ExampleGet_missing shows the zero-value + false pair returned when the
// requested type has not been stored in the context.
func ExampleGet_missing() {
	_, ok := ctxval.Get[traceID](context.Background())
	fmt.Println(ok)

	// Output: false
}

// ExampleMust_present shows that [Must] returns the stored value when found.
func ExampleMust_present() {
	ctx := ctxval.With(context.Background(), userID(7))
	fmt.Println(ctxval.Must[userID](ctx))

	// Output: 7
}

// ExampleMust_missing shows that [Must] returns the zero value of T when
// nothing has been stored — it never panics.
func ExampleMust_missing() {
	zero := ctxval.Must[userID](context.Background())
	fmt.Println(zero)

	// Output: 0
}

// ExampleOr_present shows that [Or] returns the stored value when found,
// ignoring the fallback.
func ExampleOr_present() {
	ctx := ctxval.With(context.Background(), requestID("req-42"))
	fmt.Println(ctxval.Or(ctx, requestID("fallback")))

	// Output: req-42
}

// ExampleOr_missing shows that [Or] returns the caller-supplied fallback
// when nothing has been stored.
func ExampleOr_missing() {
	fmt.Println(ctxval.Or(context.Background(), requestID("fallback")))

	// Output: fallback
}
