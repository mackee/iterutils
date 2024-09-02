package async

import (
	"context"
	"iter"
	"sync"
)

// Map applies the function `fn` to each element of the iterator `it`,
// producing a sequence of results of type U.
//
// The function `fn` is executed asynchronously for all elements, meaning
// that multiple invocations of `fn` may run concurrently on different
// elements returned by `it`.
//
// The resulting iter.Seq maintains the order of the original iterator `it`,
// ensuring that the output sequence matches the order of the input elements.
func Map[E any, U any](it iter.Seq[E], fn func(E) U) iter.Seq[U] {
	return func(yield func(U) bool) {
		ctx, cancel := context.WithCancel(context.Background())
		ch := make(chan yieldResult[U])
		wg := &sync.WaitGroup{}
		var index int
		for v := range it {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				u := fn(v)
				select {
				case ch <- yieldResult[U]{
					t:     &u,
					index: index,
				}:
				case <-ctx.Done():
				}
			}(index)
			index++
		}
		go func() {
			wg.Wait()
			close(ch)
		}()

		collectFromChannel(ctx, cancel, ch, func(t *U) bool {
			return yield(*t)
		})
	}
}

// Map2 applies the function `fn` to each pair of elements (V, W) from the iterator `it`,
// producing a sequence of pairs of results (V2, W2).
//
// The function `fn` is executed asynchronously for all pairs, meaning that multiple
// invocations of `fn` may run concurrently on different pairs returned by `it`.
//
// The resulting iter.Seq2 maintains the order of the original iterator `it`, ensuring
// that the output sequence of pairs matches the order of the input pairs.
func Map2[V any, W any, V2 any, W2 any](it iter.Seq2[V, W], fn func(V, W) (V2, W2)) iter.Seq2[V2, W2] {
	return func(yield func(V2, W2) bool) {
		ctx, cancel := context.WithCancel(context.Background())
		type result struct {
			v2 *V2
			w2 *W2
		}

		ch := make(chan yieldResult[result])
		wg := &sync.WaitGroup{}
		var index int
		for v, w := range it {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				v2, w2 := fn(v, w)
				select {
				case ch <- yieldResult[result]{
					t:     &result{v2: &v2, w2: &w2},
					index: index,
				}:
				case <-ctx.Done():
				}
			}(index)
			index++
		}
		go func() {
			wg.Wait()
			close(ch)
		}()

		collectFromChannel(ctx, cancel, ch, func(t *result) bool {
			return yield(*t.v2, *t.w2)
		})
	}
}

type yieldResult[T any] struct {
	t     *T
	index int
}

func collectFromChannel[T any](ctx context.Context, cancel func(), ch <-chan yieldResult[T], fn func(t *T) bool) {
	defer cancel()
	var currentIndex int
	m := make(map[int]*T)
	ret := make(chan *T)
	go func() {
		defer close(ret)
		for yr := range ch {
			m[yr.index] = yr.t
			for {
				if v, ok := m[currentIndex]; ok {
					select {
					case ret <- v:
					case <-ctx.Done():
						return
					}
				} else {
					break
				}
				currentIndex++
			}
		}
	}()

	for t := range ret {
		if !fn(t) {
			break
		}
	}
}
