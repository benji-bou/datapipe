package datapipe

import "github.com/benji-bou/chantools"

type Outputable[T any] interface {
	Output(input Input[T]) <-chan error
}

type Output[T any] func(input Input[T]) <-chan error

func (o Output[T]) Output(input Input[T]) <-chan error {
	return o(input)
}

func NewOutput[T any](out func(input <-chan T, err <-chan error) <-chan error) Output[T] {
	return Output[T](func(input Input[T]) <-chan error {
		ci, cie := input.Input()
		oe := out(ci, cie)
		return oe
	})
}

type OutputRaw func(input Input[[]byte]) <-chan error

func (o OutputRaw) Output(input Input[[]byte]) <-chan error {
	return o(input)
}

func Broadcast[I any](input Input[I], output ...Output[I]) <-chan error {
	inputs := Duplicate(input, len(output))
	errs := make([]<-chan error, len(output))
	for i, o := range output {
		errs[i] = o.Output(inputs[i])
	}
	return chantools.MergeChan(errs...)
}
