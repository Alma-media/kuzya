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
