package api

import (
	"fmt"
	"log"
	"regexp"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Handle func(deviceID string) (state string, err error)

// SetHandle applies an absolute state carried in the message payload
// ("ON"/"OFF") and returns the resulting state.
type SetHandle func(deviceID, requested string) (state string, err error)

func inputURL(method string) string {
	return fmt.Sprintf("/trig.in.%s", method)
}

func outputURL(method, id string) string {
	return fmt.Sprintf("/trig.out.%s/%s", method, id)
}

func pattern(method string) string {
	return `^` + regexp.QuoteMeta(method) + `\/(.*)$`
}

// RegisterHandler subscribes to /trig.in.<method>/+ and answers every
// message with the result of handle, published to /trig.out.<method>/<id>.
func RegisterHandler(client mqtt.Client, method string, handle Handle) error {
	return RegisterSetter(client, method, method, func(deviceID, _ string) (string, error) {
		return handle(deviceID)
	})
}

// RegisterSetter is RegisterHandler with two differences: the message
// payload is handed to the handler (absolute set), and the result goes to
// /trig.out.<outMethod>/<id>, which may differ from the input method. The
// "setswitch" endpoint uses outMethod "switch" so the new state travels the
// exact same elsa fan-out (relays, LEDs) and mirrors as a toggle does.
func RegisterSetter(client mqtt.Client, method, outMethod string, handle SetHandle) error {
	var (
		path  = inputURL(method)
		expr  = regexp.MustCompile(pattern(path))
		token = client.Subscribe(path+"/+", 0, func(client mqtt.Client, msg mqtt.Message) {
			if !expr.MatchString(msg.Topic()) {
				log.Println("invalid device format")

				return
			}

			deviceID := expr.FindStringSubmatch(msg.Topic())[1]

			payload, err := handle(deviceID, string(msg.Payload()))
			if err != nil {
				log.Printf("cannot resolve state for %q: %s", deviceID, err)

				return
			}

			client.Publish(outputURL(outMethod, deviceID), 0, false, payload).Wait()
		})
	)

	token.Wait()

	return token.Error()
}
