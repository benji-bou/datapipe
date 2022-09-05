package ws

import (
	"fmt"

	"github.com/benji-bou/chantools"
	"github.com/benji-bou/datapipe"
	"github.com/benji-bou/wsocket"
)

func Input(cl *wsocket.Socket) datapipe.Inputable[[]byte] {
	ci, cerr := chantools.ChanBuffErrGenerator(func(c chan<- []byte, err chan<- error) {
	L:
		for {
			select {
			case event, ok := <-cl.Read():
				if !ok {
					err <- fmt.Errorf("Websocket closed. Read channel closed")
					break L
				}
				c <- event
			case e, ok := <-cl.Error():
				if !ok {
					err <- fmt.Errorf("Websocket closed. Error channel closed")
					break L
				} else {
					err <- fmt.Errorf("Error received data from ws: %w\n", e)
					break L
				}

			}
		}
	}, 1000)
	return datapipe.NewInput(ci, cerr)
}

func Output(cl *wsocket.Socket) datapipe.Outputable[[]byte] {
	return datapipe.NewOutput(func(input <-chan []byte, err <-chan error) <-chan error {
		errWs := chantools.ChanGenerator(func(c chan<- error) {
			for i := range input {
				eSend := cl.SendMessage(i)
				if eSend != nil {
					c <- eSend
				}
			}
		})
		return chantools.MergeChan(err, errWs)
	})
}
