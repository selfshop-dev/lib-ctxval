# lib-ctxval

[![CI](https://github.com/selfshop-dev/lib-ctxval/actions/workflows/ci.yml/badge.svg)](https://github.com/selfshop-dev/lib-ctxval/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/selfshop-dev/lib-ctxval/branch/main/graph/badge.svg)](https://codecov.io/gh/selfshop-dev/lib-ctxval)
[![Go Report Card](https://goreportcard.com/badge/github.com/selfshop-dev/lib-ctxval)](https://goreportcard.com/report/github.com/selfshop-dev/lib-ctxval)
[![Go version](https://img.shields.io/github/go-mod/go-version/selfshop-dev/lib-ctxval)](go.mod)
[![License](https://img.shields.io/github/license/selfshop-dev/lib-ctxval)](LICENSE)

Типобезопасное хранение и извлечение значений в [context.Context] с использованием дженерик-ключей. Проект организации [selfshop-dev](https://github.com/selfshop-dev).

### Installation

```bash
go get -u github.com/selfshop-dev/lib-ctxval
```

## Overview

Библиотека `ctxval` решает проблему коллизий ключей при работе с `context.WithValue` из стандартной библиотеки. Вместо ручного объявления непрозрачных ключей каждый тип данных автоматически получает изолированный слот в контексте благодаря дженерикам.

```go
type RequestID string

ctx = ctxval.With(ctx, RequestID("req-42"))
id, ok := ctxval.Get[RequestID](ctx)
```

### Быстрый старт

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

## Хранение значений

Функция `With` сохраняет значение произвольного типа в контекст. Ключом выступает сам тип `T`, поэтому разные типы никогда не конфликтуют между собой.

```go
type requestID string
type traceID string

ctx := ctxval.With(context.Background(), requestID("abc"))
ctx = ctxval.With(ctx, traceID("xyz")) // безопасное добавление другого типа
```

## Извлечение значений

Библиотека предоставляет три функции для чтения значений, покрывающие основные сценарии использования.

| Функция | Возвращает | Поведение при отсутствии |
|---------|-----------|-------------------------|
| `Get[T](ctx)` | `(T, bool)` | zero value + `false` |
| `Must[T](ctx)` | `T` | zero value (без паники) |
| `Or(ctx, fallback)` | `T` | переданный `fallback` |

```go
// Get — проверка наличия значения
if id, ok := ctxval.Get[requestID](ctx); ok {
	// используем id
}

// Must — когда значение обязательно, но допустим zero value
trace := ctxval.Must[traceID](ctx)

// Or — значение или запасной вариант
token := ctxval.Or(ctx, authToken("default"))
```

## Ограничения

Каждый тип `T` занимает ровно один слот в цепочке контекстов. Повторное сохранение значения того же типа перезаписывает предыдущее, но не удаляет его — поведение идентично `context.WithValue`.

```go
ctx := ctxval.With(context.Background(), requestID("first"))
ctx = ctxval.With(ctx, requestID("second")) // "first" затенён

id, _ := ctxval.Get[requestID](ctx) // вернёт "second"
```

При хранении интерфейсных типов (например, `error`) извлекаемое значение сохраняет свой конкретный динамический тип. Для работы с ним применяются стандартные правила type assertion.

## Лицензия

[`MIT`](LICENSE) © 2026-present [`selfshop-dev`](https://github.com/selfshop-dev)