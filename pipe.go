package datapipe

import (
	"time"

	"github.com/benji-bou/chantools"
)

type Pipeable[I any, O any] interface {
	Pipe(input Input[I]) Input[O]
}

type Pipe[I any, O any] func(input Input[I]) Input[O]

func (p Pipe[I, O]) Pipe(input Input[I]) Input[O] {
	return p(input)
}

func Map[I any, O any](mapper func(input I) (O, error)) Pipe[I, O] {
	return func(input Input[I]) Input[O] {
		ci, cie := input.Input()
		mi, mie := chantools.MapChan(ci, mapper)
		return NewInput(mi, chantools.MergeChan(cie, mie))
	}
}

func TimeBuffered[I any](duration time.Duration) Pipe[I, []I] {
	return func(input Input[I]) Input[[]I] {
		ci, cie := input.Input()
		cRes := chantools.ChanGenerator(func(c chan<- []I) {
			tick := time.NewTicker(duration)
			buffer := make([]I, 0)
			for {
				select {
				case <-tick.C:
					if ci == nil && len(buffer) == 0 {
						return
					}
					c <- buffer
					buffer = buffer[:0]
				case data, ok := <-ci:
					if !ok {
						ci = nil
					}
					buffer = append(buffer, data)
				}
			}

		})
		return NewInput(cRes, cie)
	}
}

func ConcatPipe[I any]() Pipe[I, []I] {
	return func(input Input[I]) Input[[]I] {
		return Concat(input)
	}
}

func FlattenMapPipe[K comparable, V any](strat FlattenStrat) Pipe[map[K]V, map[K]V] {
	return func(input Input[map[K]V]) Input[map[K]V] {
		return FlattenMap(input, strat)
	}

}
