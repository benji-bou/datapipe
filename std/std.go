package std

import (
	"fmt"
	"os"

	"github.com/benji-bou/chantools"
	"github.com/benji-bou/datapipe"
)

func Output[T any]() datapipe.Outputable[T] {
	return datapipe.Output[T](func(input datapipe.Inputable[T]) <-chan error {
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
	})
}
