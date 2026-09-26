# go-do API guide

This guide describes the repository's current public API. Check the source or `go doc` when behavior is important and this reference may be stale. The module path is `go.osspkg.com/do`; `go.mod` declares Go 1.26.0.

## Choose a package

| Package | Use it for |
| --- | --- |
| `go.osspkg.com/do` | Generic slice and map operations, simple numeric helpers, conditional expressions, panic recovery, and asynchronous calls. |
| `go.osspkg.com/do/mapreduce` | Concurrent mapping followed by reduction of completed results. |
| `go.osspkg.com/do/workerpool` | Reusable bounded workers that accept identified tasks and publish results on a channel. |
| `go.osspkg.com/do/pipline` | Context-aware handlers connected by comparable state values. Preserve the existing spelling `pipline`. |
| `go.osspkg.com/do/monad` | A small generic value-plus-error `Result` with `Bind`. |
| `go.osspkg.com/do/words` | Split strings or byte slices into tokens using default or custom rune detectors. |

## Root package

### Slices

- Iteration and conversion: `Each`, `Convert`, `Reduce`, `Filter`, `Treat`, `TreatValue`, `Entries`, `ToMap`.
- Composition and grouping: `Join`, `Chunk`, `Diff`, `Unique`, `Exclude`.
- Search: `IndexOf`, `LastIndexOf`, `Include`.
- Mutation: `Push`, `Pop`, `Shift`, `Unshift`, `Reverse`, `Splice`. These mutate the supplied slice or slice pointer. `Copy` returns a new slice; its `start` and `end` indexes are inclusive and clamped to the input bounds.

`Reduce` starts with the element type's zero value; it does not accept a separate initial accumulator. Empty inputs return zero or an empty result according to the function.

### Maps

- Traversal and transformation: `EachMap`, `ConvertMap`, `FilterMap`, `TreatMap`, `TreatMapValue`, `ReduceMap`, `ToSlice`.
- Key/value operations: `Keys`, `Values`, `JoinMap`, `FlipMap`, `DivideMap`, `CombineMap`.

`Keys`, `Values`, `DivideMap`, `ReduceMap`, and `ToSlice` use sorted keys, constrained by the library's `Comparable` type set. Other helpers that range directly over a map do not promise iteration order. Key collisions in conversion, joining, flipping, or combining use the last value assigned during iteration; where input map iteration is involved, which colliding entry wins may be nondeterministic.

### Numbers and conditionals

- `MinMax` and `MinMaxTime` clamp a value to a minimum and maximum.
- `Range` generates an ascending inclusive sequence using a positive step. It returns an empty slice for non-positive steps and stops if adding the step does not increase the current value.
- `Sum` adds values supported by `Comparable`; `Average` supports `Summable` numeric types and uses the type's integer or floating-point division semantics. Calling `Average` with no values divides by zero for integer types.
- `If`, `IfElse`, and their `Func` variants provide conditional selection. Function arguments are useful for lazy evaluation. `DoIf` supports an `ElseIf`/`Else` chain and lazy `ElseIfFunc`/`ElseFunc` branches.

`Summable` and `Comparable` are generic type constraints intended for these APIs, not runtime validators.

### Panic and asynchronous helpers

- `Recovery` runs a function and returns a recovered panic as an error; it does not recover panics in other goroutines.
- `Try` runs optional try/catch/finally callbacks. Panics raised by callbacks are recovered by its internal wrappers.
- `Trace` formats runtime call frames.
- `Async` starts one function in a goroutine and optionally reports a panic through the error callback. It does not provide a join handle.
- `AsyncGroup` starts one goroutine per supplied function, waits for all of them, and returns their errors. It passes the supplied context but does not itself cancel sibling functions after one fails; use a bounded worker pool for large workloads.

## `mapreduce`

`mapreduce.New(ctx, items, mapper, reducer, initial, workers)` allocates a result channel sized to `len(items)`, runs mapper calls with an `errgroup` concurrency limit, then reduces results as they arrive. Large input slices therefore also require a potentially large result buffer. `workers` must be greater than zero. With multiple workers, completion and reduction order are nondeterministic; reducers that require input order should use one worker or an order-preserving design. Mapper functions should observe context cancellation. A reducer error cancels the mapper context and waits for mapper goroutines to return.

## `workerpool`

- `New(workers, handler)` creates a pool; non-positive worker counts are clamped to one.
- `Start` starts workers once; repeated calls do nothing.
- `Send` returns whether the task was accepted before shutdown. It returns `false` after closure or cancellation.
- `Receive` exposes a result channel buffered to the worker count.
- `Close` cancels workers, waits for senders and workers, and closes the result channel. It is safe to call repeatedly.

The handler receives the pool context. `Close` waits for active handlers, so handlers that need responsive shutdown should return when the context is canceled. Continue receiving results during normal work because the result buffer is bounded; cancellation unblocks workers that cannot publish another result.

## `pipline`

Create a `Pipline[S, T]` with `New`, register `Pipe[S, T]` transitions using `Set`, and execute using `Do(ctx, state, value)`. `Set` rejects a nil handler, duplicate current states, and transitions whose current and next states are equal. `Do` follows transitions until a state has no registered handler, the context is canceled, or a handler returns an error.

## `monad`

`Some(value)` creates a successful `Result[T]`. `Bind(result, next)` skips `next` if the result already holds an error; otherwise it stores the returned value or error. `Result.Return()` retrieves the pair. The contained fields are private, so construct values with `Some` and `Bind`.

## `words`

`Strings(s)` and `Bytes(b)` use the default detectors. `NewString` and `NewBytes` return a `Words` interface that can be configured with `SetBlock`, `SetDigital`, and `SetSymbol`, or reset with the matching `UseDefault...` methods. `Words.Strings()` and `Words.Bytes()` tokenize from the beginning on each call. The default block detector groups letters and digits; digit and symbol detectors control how adjacent tokens are grouped. Tokenization uses `bufio.Scanner` with its default token limit and ignores scanner errors, so very long tokens can cause incomplete output without an error being returned.
