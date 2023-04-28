package std

import (
	"fmt"
	"os"

	"github.com/benji-bou/chantools"
	"github.com/benji-bou/datapipe"
	"github.com/k0kubun/pp"
)

func Output[T any](input datapipe.Input[T]) <-chan error {

	return chantools.ChanGenerator(func(c chan<- error) {
		ci, cie := input.Input()
		for {
			select {
			case data, ok := <-ci:
				if !ok {
					return
				}
				fmt.Println(data)
			case err, ok := <-cie:
				if !ok {
					return
				}
				fmt.Fprintln(os.Stderr, err)
			}
		}
	})
}

func PrettyOutput[T any](input datapipe.Input[T]) <-chan error {

	return chantools.ChanGenerator(func(c chan<- error) {
		ci, cie := input.Input()
		for {
			select {
			case data, ok := <-ci:
				if !ok {
					return
				}
				pp.Println(data)
			case err, ok := <-cie:
				if !ok {
					return
				}
				fmt.Fprintln(os.Stderr, err)
			}
		}
	})
}
