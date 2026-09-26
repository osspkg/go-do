# go-do

[![Go version](https://img.shields.io/github/go-mod/go-version/osspkg/go-do)](https://go.dev/doc/install)
[![CI](https://github.com/osspkg/go-do/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/osspkg/go-do/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/go.osspkg.com/do.svg)](https://pkg.go.dev/go.osspkg.com/do)
[![License](https://img.shields.io/github/license/osspkg/go-do)](LICENSE)

`go-do` is a Go utility library with generic helpers for slices and maps, numeric operations, conditional expressions, panic recovery, and asynchronous work. It also provides focused packages for worker pools, map-reduce, state pipelines, result handling, and word tokenization.

## Requirements

- Go 1.26.0 or newer, as declared in [`go.mod`](go.mod).

## Installation

```sh
go get go.osspkg.com/do
```

Import the root package and any subpackage you need:

```go
import (
    "go.osspkg.com/do"
    "go.osspkg.com/do/mapreduce"
    "go.osspkg.com/do/workerpool"
)
```

## Quick start

```go
package main

import (
    "fmt"

    "go.osspkg.com/do"
)

func main() {
    values := do.Filter([]int{1, 2, 3, 4, 5}, func(value, index int) bool {
        return value%2 == 1
    })
    labels := do.Convert(values, func(value, index int) string {
        return fmt.Sprintf("item-%d", value)
    })

    fmt.Println(labels) // [item-1 item-3 item-5]
}
```

## Packages and features

### Root package: `go.osspkg.com/do`

| Area | Functions and types | Notes |
| --- | --- | --- |
| Slices | `Each`, `Convert`, `Join`, `Chunk`, `Entries`, `Reduce`, `Filter`, `Treat`, `TreatValue`, `Diff`, `Unique`, `IndexOf`, `LastIndexOf`, `Include`, `Exclude`, `Copy`, `Splice`, `Pop`, `Push`, `Shift`, `Unshift`, `Reverse`, `ToMap` | `Splice`, `Push`, `Unshift`, `Pop`, `Shift`, and `Reverse` mutate the supplied slice or slice pointer. |
| Maps | `EachMap`, `ConvertMap`, `FilterMap`, `TreatMap`, `TreatMapValue`, `Keys`, `Values`, `JoinMap`, `FlipMap`, `DivideMap`, `CombineMap`, `ReduceMap`, `ToSlice` | `Keys`, `Values`, `DivideMap`, `ReduceMap`, and `ToSlice` use sorted keys. Other map transformations may follow Go's unspecified map iteration order. |
| Numeric helpers | `MinMax`, `MinMaxTime`, `Range`, `Sum`, `Average` | `Range` includes both endpoints when reached and requires a positive, progressing step; otherwise it returns an empty slice or stops when the value cannot advance. |
| Conditional helpers | `If`, `IfFunc`, `IfElse`, `IfElseFunc`, `DoIf` | `DoIf.ElseIf`, `Else`, `ElseIfFunc`, and `ElseFunc` form a conditional chain; function variants evaluate only the selected branch. |
| Panic and async helpers | `Recovery`, `Trace`, `Try`, `Async`, `AsyncGroup` | `AsyncGroup` waits for all supplied functions and returns their errors. |

The `Summable` and `Comparable` type constraints define the numeric and ordered types accepted by the generic helpers.

### `go.osspkg.com/do/workerpool`

A bounded worker pool with generic task and result types.

| API | Use |
| --- | --- |
| `Task[I, T]`, `Result[I, T]` | Carry a task ID and input data, or a task ID, result value, and error. |
| `Pool[I, T, R]` | Own the task and result channels and worker lifecycle. |
| `New[I, T, R]` | Create a pool with a fixed worker count and task handler. |
| `Pool.Start` | Start the workers; repeated calls have no effect. |
| `Pool.Send` | Submit a task; returns `false` after shutdown or when the pool is closing. |
| `Pool.Receive` | Get the result channel. |
| `Pool.Close` | Cancel workers and close the result channel; safe to call more than once. |

Read results while work is running to avoid filling the bounded result buffer during normal processing.

### `go.osspkg.com/do/mapreduce`

| API | Use |
| --- | --- |
| `New[T, R, O]` | Map items concurrently and pass completed results to a reducer. |

Set `workers` to a positive number. Mapping and reducing order is nondeterministic when more than one worker is used, so reducers should not depend on input order. A reducer error cancels the mapping context and waits for mapper goroutines to return; mapper functions should honor their context.

### `go.osspkg.com/do/pipline`

| API | Use |
| --- | --- |
| `Pipe[S, T]` | Define a state transition and its context-aware handler. |
| `Pipline[S, T]` | Store transitions keyed by their current state. |
| `New[S, T]` | Create an empty pipeline. |
| `Pipline.Set` | Register a transition; returns an error for a nil handler, duplicate current state, or identical current and next states. |
| `Pipline.Do` | Run handlers from a starting state until no transition exists, the context is canceled, or a handler returns an error. |

### `go.osspkg.com/do/monad`

| API | Use |
| --- | --- |
| `Result[T]` | Hold a value and an error for chained operations. |
| `Some[T]` | Create a successful result from a value. |
| `Bind[T, R]` | Apply a function to a successful result and propagate errors. |
| `Result.Return` | Retrieve the result's value and error. |

### `go.osspkg.com/do/words`

| API | Use |
| --- | --- |
| `Words` | Interface for tokenizing to strings or byte slices and configuring token categories. |
| `NewString`, `NewBytes` | Create a tokenizer with default rune detectors. |
| `Strings`, `Bytes` | Tokenize a string or byte slice with the default block and symbol detectors. |
| `Words.Strings`, `Words.Bytes` | Tokenize using the instance's current detector configuration. |
| `UseDefaultBlock`, `SetBlock` | Select default or custom runes that form word blocks. |
| `UseDefaultDigital`, `SetDigital` | Select default or custom digit runes. |
| `UseDefaultSymbol`, `SetSymbol` | Select default or custom symbol runes. |

## Testing

Run the complete test suite from the repository root:

```sh
go test ./...
```

The project Makefile also provides `make tests`, `make lint`, `make build`, and `make ci`. `make ci` installs the latest `goppy` tool and runs setup, license, lint, test, and build targets; see the [Makefile](Makefile) before running it locally.

## Contributing

Pull requests are welcome. Include tests for behavior changes and run `go test ./...` before submitting. Keep public API and import-path compatibility in mind.

## License

This project is licensed under the [BSD 3-Clause License](LICENSE).
