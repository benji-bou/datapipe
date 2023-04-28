package mqtt

import (
	"context"
	"errors"

	"github.com/benji-bou/chantools"
	"github.com/benji-bou/datapipe"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func Input(ctx context.Context, topic string, qos byte, client mqtt.Client) datapipe.Input[[]byte] {
	ci, cerr := chantools.ChanBuffErrGenerator(func(c chan<- []byte, err chan<- error) {
		token := client.Subscribe(topic, qos, func(client mqtt.Client, m mqtt.Message) {
			m.Ack()
			data := m.Payload()
			c <- data
		})
		if token.Wait() && token.Error() != nil {
			err <- token.Error()
		}
		<-ctx.Done()
	}, 255)
	return datapipe.NewInput(ci, cerr)
}

func Output(topic string, client mqtt.Client) datapipe.Output[[]byte] {
	return datapipe.NewOutput(func(input <-chan []byte, err <-chan error) <-chan error {
		errMqtt := chantools.ChanGenerator(func(c chan<- error) {
			for i := range input {
				if !client.IsConnected() {
					c <- errors.New("client has been disconnected")
					return
				}
				token := client.Publish(topic, 1, false, i)
				w := token.Wait()
				if w && token.Error() != nil {
					c <- token.Error()
				} else if !w {
					c <- errors.New("wait is false")
				}
			}
		})
		return chantools.MergeChan(err, errMqtt)
	})
}
