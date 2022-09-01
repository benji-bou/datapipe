package datapipe

import "github.com/benji-bou/chantools"

type Inputable[T any] interface {
	Input() (<-chan T, <-chan error)
}
type Input[T any] func() (<-chan T, <-chan error)

func (in Input[T]) Input() (<-chan T, <-chan error) {
	return in()
}

func NewInput[T any](i <-chan T, e <-chan error) Inputable[T] {
	return Input[T](func() (<-chan T, <-chan error) { return i, e })
}

type InputRaw func() (<-chan []byte, <-chan error)

func (ir InputRaw) Input() (<-chan []byte, <-chan error) {
	return ir()
}

func Duplicate[T any](i Inputable[T], count int) []Inputable[T] {
	res := make([]Inputable[T], count)
	ci, cie := i.Input()
	bci := chantools.NewBrokerChan(ci)
	bcie := chantools.NewBrokerChan(cie)
	for i := 0; i < count; i++ {
		rci := bci.Subscribe()
		rcie := bcie.Subscribe()
		res[i] = NewInput(rci, rcie)
	}
	return res
}

func Merge[T any](input ...Inputable[T]) Inputable[T] {
	cis := make([]<-chan T, len(input))
	ces := make([]<-chan error, len(input))
	for i, in := range input {
		inC, inCe := in.Input()
		cis[i] = inC
		ces[i] = inCe
	}
	return NewInput(chantools.MergeChan(cis...), chantools.MergeChan(ces...))
}
