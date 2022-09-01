package datapipe

import (
	"time"

	"github.com/benji-bou/chantools"
)

type Pipeable[I any, O any] interface {
	Pipe(input Inputable[I]) Inputable[O]
}

type Pipe[I any, O any] func(input Inputable[I]) Inputable[O]

func (p Pipe[I, O]) Pipe(input Inputable[I]) Inputable[O] {
	return p(input)
}

func Map[I any, O any](mapper func(input I) (O, error)) Pipeable[I, O] {
	return Pipe[I, O](func(input Inputable[I]) Inputable[O] {
		ci, cie := input.Input()
		mi, mie := chantools.MapChan(ci, mapper)
		return NewInput(mi, chantools.MergeChan(cie, mie))
	})
}

func TimeBuffered[I any](duration time.Duration) Pipeable[I, []I] {
	return Pipe[I, []I](func(input Inputable[I]) Inputable[[]I] {
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
	})
}
