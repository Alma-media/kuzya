package api

import (
	"log"
	"regexp"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var expr = regexp.MustCompile(`^\/trig-in\/(.*)$`)

func CreateStateHandler(trig func(string) (string, error)) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		if !expr.MatchString(msg.Topic()) {
			log.Println("invalid device format")

			return
		}

		deviceID := expr.FindStringSubmatch(msg.Topic())[1]

		state, err := trig(deviceID)
		if err != nil {
			log.Printf("cannot retrieve current state: %s", err)

			return
		}

		client.Publish("/trig-out/"+deviceID, 0, false, state).Wait()
	}
}
