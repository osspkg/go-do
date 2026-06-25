package mapreduce

import (
	"context"

	"golang.org/x/sync/errgroup"
)

func New[T, R, O any](
	ctx context.Context,
	items []T,
	mapper func(context.Context, T) (R, error),
	reducer func(context.Context, O, R) (O, error),
	initial O,
	workers int,
) (O, error) {
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(workers)

	results := make(chan R, len(items))

	for _, item := range items {
		item := item
		g.Go(func() error {
			r, err := mapper(ctx, item)
			if err != nil {
				return err
			}
			select {
			case results <- r:
			case <-ctx.Done():
				return ctx.Err()
			}
			return nil
		})
	}

	go func() {
		_ = g.Wait()
		close(results)
	}()

	acc := initial
	for r := range results {
		var err error
		acc, err = reducer(ctx, acc, r)
		if err != nil {
			return acc, err
		}
	}

	if err := g.Wait(); err != nil {
		return acc, err
	}
	return acc, nil
}
