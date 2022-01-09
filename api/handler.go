package api

import (
	"fmt"
	"log"
	"regexp"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Handle func(deviceID string) (state string, err error)

func inputURL(method string) string {
	return fmt.Sprintf("/trig.in.%s", method)
}

func outputURL(method, id string) string {
	return fmt.Sprintf("/trig.out.%s/%s", method, id)
}

func pattern(method string) string {
	return `^` + regexp.QuoteMeta(method) + `\/(.*)$`
}

func RegisterHandler(client mqtt.Client, method string, handle Handle) error {
	var (
		path  = inputURL(method)
		expr  = regexp.MustCompile(pattern(path))
		token = client.Subscribe(path+"/+", 0, func(client mqtt.Client, msg mqtt.Message) {
			if !expr.MatchString(msg.Topic()) {
				log.Println("invalid device format")

				return
			}

			deviceID := expr.FindStringSubmatch(msg.Topic())[1]

			payload, err := handle(deviceID)
			if err != nil {
				log.Printf("cannot retrieve current state: %s", err)

				return
			}

			client.Publish(outputURL(method, deviceID), 0, false, payload).Wait()
		})
	)

	token.Wait()

	return token.Error()
}
