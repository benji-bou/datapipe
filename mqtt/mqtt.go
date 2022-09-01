package mqtt

import (
	"github.com/benji-bou/chantools"
	"github.com/benji-bou/datapipe"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func Input(topic string, qos byte, client mqtt.Client) datapipe.Inputable[[]byte] {
	ci, cerr := chantools.ChanBuffErrGenerator(func(c chan<- []byte, err chan<- error) {
		token := client.Subscribe(topic, qos, func(client mqtt.Client, m mqtt.Message) {
			m.Ack()
			data := m.Payload()
			c <- data
		})
		if token.Wait() && token.Error() != nil {
			err <- token.Error()
		}
	}, 255)
	return datapipe.NewInput(ci, cerr)
}

func Output(topic string, client mqtt.Client) datapipe.Outputable[[]byte] {
	return datapipe.NewOutput(func(input <-chan []byte, err <-chan error) <-chan error {
		errMqtt := chantools.ChanGenerator(func(c chan<- error) {
			for i := range input {
				token := client.Publish(topic, 1, false, i)
				if token.Wait() && token.Error() != nil {
					c <- token.Error()
				}
			}
		})
		return chantools.MergeChan(err, errMqtt)
	})
}
