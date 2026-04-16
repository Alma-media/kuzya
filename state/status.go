package state

const (
	ON  = "ON"
	OFF = "OFF"
)

func Status(status bool) string {
	if status {
		return ON
	}

	return OFF
}

type Device struct {
	ID    string `json:"device_id"`
	State string `json:"state"`
}
