/*
 *  Copyright (c) 2024-2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package mapreduce

import (
	"context"
	"errors"

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
	if workers <= 0 {
		return initial, errors.New("workers must be greater than zero")
	}

	workCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	g, ctx := errgroup.WithContext(workCtx)
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
			cancel()
			_ = g.Wait()
			return acc, err
		}
	}

	if err := g.Wait(); err != nil {
		return acc, err
	}
	return acc, nil
}
