package datapipe

import "github.com/benji-bou/chantools"

type Outputable[T any] interface {
	Output(input Inputable[T]) <-chan error
}

type Output[T any] func(input Inputable[T]) <-chan error

func (o Output[T]) Output(input Inputable[T]) <-chan error {
	return o(input)
}

func NewOutput[T any](out func(input <-chan T, err <-chan error) <-chan error) Outputable[T] {
	return Output[T](func(input Inputable[T]) <-chan error {
		ci, cie := input.Input()
		oe := out(ci, cie)
		return oe
	})
}

type OutputRaw func(input Inputable[[]byte]) <-chan error

func (o OutputRaw) Output(input Inputable[[]byte]) <-chan error {
	return o(input)
}

func Broadcast[I any](input Inputable[I], output ...Outputable[I]) <-chan error {
	inputs := Duplicate(input, len(output))
	errs := make([]<-chan error, len(output))
	for i, o := range output {
		errs[i] = o.Output(inputs[i])
	}
	return chantools.MergeChan(errs...)
}
