---
name: go-do-usage
description: Choose and use APIs from the go.osspkg.com/do Go utility library. Use when writing, reviewing, or explaining Go code that imports this module; skip for Go work that does not use go-do.
metadata:
  short-description: Use the go-do Go library
---

# Using go-do

Use the library's generic helpers and focused subpackages when they match the caller's needs. Preserve their actual mutation, ordering, error, and lifecycle behavior; do not infer semantics from familiar function names alone.

## Workflow

1. Confirm the importing module uses `go.osspkg.com/do` and inspect its `go.mod` for the required Go version.
2. Read [the API guide](references/api-guide.md) to choose a package and understand behavioral constraints.
3. Read [the examples](references/examples.md) for the relevant operation, then verify details against the current implementation or `go doc` if code and reference may have drifted.
4. Add or update tests for behavior changes using the consuming repository's conventions.

## Important usage constraints

- Import subpackages by their exact module paths. The state pipeline package is spelled `pipline` in both its package name and import path.
- Root helpers generally return new slices or maps, but slice operations such as `Reverse`, `Push`, `Unshift`, `Pop`, `Shift`, and `Splice` mutate their input. Check the API guide before assuming ownership or aliasing behavior.
- Go map iteration is unordered. Helpers that sort keys provide stable key order; do not assume that concurrent map-reduce preserves input order.
- Use a positive worker count for `mapreduce.New`, and make mapper functions honor their context. A reducer error cancels that context and waits for mapper functions to return.
- For `workerpool`, start the pool once, handle every expected result while tasks are running, and call `Close` to cancel workers and close the result channel. A handler that needs prompt shutdown must honor its context.
- `do.Range` is ascending: use a positive step that advances the value. Non-positive steps produce an empty slice; integer overflow or floating-point precision that prevents progress stops the range.

## References

- [API guide](references/api-guide.md): package map, selection guidance, and behavior notes.
- [Examples](references/examples.md): focused snippets for collection helpers, map-reduce, worker pools, pipelines, words, and result chaining.
