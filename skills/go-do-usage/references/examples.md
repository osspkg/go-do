# go-do examples

The examples use the module imports from `go.osspkg.com/do`. Each code block is a function intended to be copied into a Go package together with the listed imports.

## Filter and convert a slice

```go
func filterLabels() []string {
    values := do.Filter([]int{1, 2, 3, 4, 5}, func(value, index int) bool {
        return value%2 == 1
    })
    return do.Convert(values, func(value, index int) string {
        return fmt.Sprintf("item-%d", value)
    })
}
```

Imports: `fmt` and `go.osspkg.com/do`.

## Transform a map with stable key order

```go
func formatScores() []string {
    scores := map[string]int{"Ada": 10, "Lin": 8}
    return do.ToSlice(scores, func(score int, name string) string {
        return fmt.Sprintf("%s=%d", name, score)
    })
}
```

Imports: `fmt` and `go.osspkg.com/do`. `ToSlice` visits map entries in sorted key order.

## Run map-reduce work

```go
func sumSquares() (int, error) {
    ctx := context.Background()
    total, err := mapreduce.New(
        ctx,
        []int{1, 2, 3, 4},
        func(ctx context.Context, value int) (int, error) {
            return value * value, nil
        },
        func(ctx context.Context, total, value int) (int, error) {
            return total + value, nil
        },
        0,
        4,
    )
    if err != nil {
        return 0, err
    }
    return total, nil // 30; completion order is unspecified with multiple workers.
}
```

Imports: `context` and `go.osspkg.com/do/mapreduce`.

## Use a worker pool

```go
func runPool() error {
    pool := workerpool.New(2, func(ctx context.Context, task workerpool.Task[int, int]) (int, error) {
        return task.Data * 2, nil
    })
    pool.Start()
    defer pool.Close()

    tasks := []int{10, 20, 30}
    for id, value := range tasks {
        if !pool.Send(workerpool.Task[int, int]{ID: id, Data: value}) {
            return errors.New("worker pool is closed")
        }
    }

    for range tasks {
        result := <-pool.Receive()
        if result.Err != nil {
            return result.Err
        }
        fmt.Println(result.TaskID, result.Result)
    }
    return nil
}
```

Imports: `context`, `errors`, `fmt`, and `go.osspkg.com/do/workerpool`. Read as many results as tasks submitted; result arrival order is not guaranteed.

## Chain states through a pipeline

```go
func runPipeline() (string, error) {
    const (
        start = iota
        done
    )

    pipe := pipline.New[int, string]()
    err := pipe.Set(pipline.Pipe[int, string]{
        Current: start,
        Next:    done,
        Handler: func(ctx context.Context, value string) (string, error) {
            return value + "!", nil
        },
    })
    if err != nil {
        return "", err
    }

    state, value, err := pipe.Do(context.Background(), start, "hello")
    if err != nil {
        return "", err
    }
    if state != done {
        return "", errors.New("pipeline stopped before the done state")
    }
    return value, nil // "hello!"
}
```

Imports: `context`, `errors`, and `go.osspkg.com/do/pipline`.

## Configure word tokenization

```go
func tokenize() []string {
    tokenizer := words.NewString("Hello, go-do!")
    tokenizer.SetBlock(unicode.IsLetter)
    return tokenizer.Strings()
}
```

Imports: `unicode` and `go.osspkg.com/do/words`. The interface also supports custom digit and symbol detectors.

## Chain a value and error

```go
func parseNumber() (int, error) {
    result := monad.Some("42")
    parsed := monad.Bind(result, strconv.Atoi)
    return parsed.Return()
}
```

Imports: `strconv` and `go.osspkg.com/do/monad`.
