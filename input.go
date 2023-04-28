package datapipe

import (
	"log"

	"github.com/benji-bou/chantools"
)

type Inputable[T any] interface {
	Input() (<-chan T, <-chan error)
}
type Input[T any] func() (<-chan T, <-chan error)

func (in Input[T]) Input() (<-chan T, <-chan error) {
	return in()
}

func NewInput[T any](i <-chan T, e <-chan error) Input[T] {
	return func() (<-chan T, <-chan error) { return i, e }
}

type InputRaw func() (<-chan []byte, <-chan error)

func (ir InputRaw) Input() (<-chan []byte, <-chan error) {
	return ir()
}

func Duplicate[T any](i Input[T], count int) []Input[T] {
	res := make([]Input[T], count)
	ci, cie := i()
	bci := chantools.NewBrokerChan(ci)
	bcie := chantools.NewBrokerChan(cie)
	for i := 0; i < count; i++ {
		rci := bci.Subscribe()
		rcie := bcie.Subscribe()
		res[i] = NewInput(rci, rcie)
	}
	return res
}

func Merge[T any](input ...Input[T]) Input[T] {
	cis := make([]<-chan T, len(input))
	ces := make([]<-chan error, len(input))
	for i, in := range input {
		inC, inCe := in()
		cis[i] = inC
		ces[i] = inCe
	}
	return NewInput(chantools.MergeChan(cis...), chantools.MergeChan(ces...))
}

func First[I any](input Input[I]) Input[I] {
	ci, cie := input()
	cires := chantools.ChanGenerator[I](func(c chan<- I) {
		data := <-ci
		c <- data
	})
	return NewInput(cires, cie)
}

func Sequence[I any](input Input[[]I]) Input[I] {
	ci, cie := input()
	cires := chantools.ChanGenerator(func(c chan<- I) {

		for arInput := range ci {
			for _, in := range arInput {
				c <- in
			}
		}
	})
	return NewInput(cires, cie)
}

func MapInput[I any, O any](input Input[I], mapper func(input I) (O, error)) Input[O] {
	return Map(mapper)(input)
}

func Concat[I any](input Input[I]) Input[[]I] {
	ci, cie := input()
	cires := chantools.ChanGenerator(func(c chan<- []I) {
		res := make([]I, 0)
		for in := range ci {
			res = append(res, in)
		}
		c <- res
	})
	return NewInput(cires, cie)
}

func Scan[I any, O any](input Input[I], initialValue O, scanner func(elem I, res O) (O, error)) Input[O] {
	ci, cie := input()
	ciRes, cieRes := chantools.ChanErrGenerator(func(c chan<- O, e chan<- error) {
		res := initialValue
		var err error
		for i := range ci {
			res, err = scanner(i, res)
			if err != nil {
				e <- err
				err = nil
			} else {
				c <- res
			}
		}
	})
	return NewInput(ciRes, chantools.MergeChan(cie, cieRes))
}

type FlattenStrat int

const (
	ScanStrat   FlattenStrat = 0
	ConcatStrat FlattenStrat = 1
)

func FlattenMap[K comparable, V any](input Input[map[K]V], strat FlattenStrat) Input[map[K]V] {
	switch strat {
	case ScanStrat:
		return flattenMapScan(input)
	case ConcatStrat:
		return flattenMapConcat(input)
	}
	return nil
}

func flattenMapScan[K comparable, V any](input Input[map[K]V]) Input[map[K]V] {
	return Scan(input, make(map[K]V), func(elem map[K]V, res map[K]V) (map[K]V, error) {
		m, err := chantools.MergeMap(res, elem)
		if err != nil {
			log.Printf("a flatten operation with strat Scan failed: %v\n", err)
		}
		return m, err
	})
}

func flattenMapConcat[K comparable, V any](input Input[map[K]V]) Input[map[K]V] {
	return Map(func(all []map[K]V) (map[K]V, error) {
		return chantools.MergeMap(all...)
	})(Concat(input))
}
