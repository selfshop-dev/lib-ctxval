# lib-ctxval

[![CI](https://github.com/selfshop-dev/lib-ctxval/actions/workflows/ci.yml/badge.svg)](https://github.com/selfshop-dev/lib-ctxval/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/selfshop-dev/lib-ctxval/branch/main/graph/badge.svg)](https://codecov.io/gh/selfshop-dev/lib-ctxval)
[![Go Report Card](https://goreportcard.com/badge/github.com/selfshop-dev/lib-ctxval)](https://goreportcard.com/report/github.com/selfshop-dev/lib-ctxval)
[![Go version](https://img.shields.io/github/go-mod/go-version/selfshop-dev/lib-ctxval)](go.mod)
[![License](https://img.shields.io/github/license/selfshop-dev/lib-ctxval)](LICENSE)

Type-safe storage and retrieval of values in [context.Context] using generic keys. A project by [selfshop-dev](https://github.com/selfshop-dev).

### Installation

```bash
go get -u github.com/selfshop-dev/lib-ctxval
```

## Overview

`ctxval` solves the key collision problem when working with `context.WithValue` from the standard library. Instead of manually declaring opaque keys, each data type automatically gets an isolated slot in the context thanks to generics.

```go
type RequestID string

ctx = ctxval.With(ctx, RequestID("req-42"))
id, ok := ctxval.Get[RequestID](ctx)
```

### Quick Start

```go
package main

import (
	"context"
	"fmt"

	ctxval "github.com/selfshop-dev/lib-ctxval"
)

type userID int64

func main() {
	ctx := ctxval.With(context.Background(), userID(123))

	id, ok := ctxval.Get[userID](ctx)
	if ok {
		fmt.Println("User ID:", id)
	}
}
```

## Storing Values

`With` stores a value of any type in the context. The type `T` itself serves as the key, so different types never conflict with each other.

```go
type requestID string
type traceID string

ctx := ctxval.With(context.Background(), requestID("abc"))
ctx = ctxval.With(ctx, traceID("xyz")) // safely adding a different type
```

## Retrieving Values

The library provides three functions for reading values, covering the most common use cases.

| Function | Returns | When absent |
|---------|---------|-------------|
| `Get[T](ctx)` | `(T, bool)` | zero value + `false` |
| `Must[T](ctx)` | `T` | zero value (no panic) |
| `Or(ctx, fallback)` | `T` | provided `fallback` |

```go
// Get — check whether a value is present
if id, ok := ctxval.Get[requestID](ctx); ok {
	// use id
}

// Must — when the value is expected but a zero value is acceptable
trace := ctxval.Must[traceID](ctx)

// Or — value or a fallback
token := ctxval.Or(ctx, authToken("default"))
```

## Limitations

Each type `T` occupies exactly one slot in the context chain. Storing a value of the same type again shadows the previous one but does not remove it — the behavior is identical to `context.WithValue`.

```go
ctx := ctxval.With(context.Background(), requestID("first"))
ctx = ctxval.With(ctx, requestID("second")) // "first" is shadowed

id, _ := ctxval.Get[requestID](ctx) // returns "second"
```

When storing interface types (e.g. `error`), the retrieved value retains its concrete dynamic type. Standard type assertion rules apply when working with it.

## License

[`MIT`](LICENSE) © 2026-present [`selfshop-dev`](https://github.com/selfshop-dev)