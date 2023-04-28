package file

import (
	"bufio"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	"github.com/benji-bou/chantools"
	"github.com/benji-bou/datapipe"
)

func Input(path string, chunkLength int) datapipe.Input[[]byte] {
	return datapipe.Input[[]byte](func() (<-chan []byte, <-chan error) {

		return chantools.ChanBuffErrGenerator(func(data chan<- []byte, err chan<- error) {
			f, e := os.Open(path)
			if e != nil {
				err <- e
				return
			}
			defer f.Close()
			r := bufio.NewReader(f)
			for {
				buf := make([]byte, chunkLength) //the chunk size
				n, e := r.Read(buf)              //loading chunk into buffer
				buf = buf[:n]
				if n == 0 {
					if e == io.EOF {
						break
					}
					if e != nil {
						err <- e
						break
					}
					return
				}
				data <- buf
			}
		}, 1000)
	})
}

func DirInput(path string, regexFilePath string, chunkLength int) datapipe.Input[[]byte] {
	return datapipe.Input[[]byte](func() (<-chan []byte, <-chan error) {
		return chantools.ChanErrGenerator(func(c chan<- []byte, e chan<- error) {
			err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
				if info.Mode().IsRegular() == true && info.IsDir() == false {
					m, err := regexp.MatchString(regexFilePath, path)
					if err != nil {
						return err
					}
					if m {
						cData, cErr := Input(path, chunkLength).Input()
						for {
							select {
							case data, ok := <-cData:
								if !ok {
									return nil
								}
								c <- data
							case err, ok := <-cErr:
								if !ok {
									return nil
								}
								return err
							}
						}
					}
				}
				return nil
			})
			if err != nil {
				e <- err
			}
		})
	})
}
