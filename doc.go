// Package ctxval provides type-safe storage and retrieval of values in a
// [context.Context] using generic key types.
//
// The standard library's [context.WithValue] requires an opaque key to avoid
// collisions between packages. ctxval encodes the key as a zero-size generic
// struct, so every distinct type T gets its own collision-free slot without
// any manual key declaration:
//
//	type RequestID string
//
//	ctx = ctxval.With(ctx, RequestID("req-42"))
//
//	id, ok := ctxval.Get[RequestID](ctx)
//	// id == "req-42", ok == true
//
// # Retrieving values
//
// Three retrieval functions cover the common patterns:
//
//   - [Get] returns the value and a boolean, mirroring the map-lookup idiom.
//   - [Must] returns the value or the zero value of T; never panics.
//   - [Or] returns the value or a caller-supplied fallback.
//
// # Limitations
//
// Each type T occupies exactly one slot per context chain. Storing a second
// value of the same type shadows — but does not remove — the previous one,
// which is identical behaviour to [context.WithValue].
//
// Interface types (e.g. T = error) can be stored, but the retrieved value
// has the concrete dynamic type that was passed to [With]; callers must
// account for this with the usual type-assertion rules.
//
// # Concurrency
//
// All functions in this package are safe for concurrent use. They rely
// exclusively on [context.Context], which is immutable after creation.
package ctxval
