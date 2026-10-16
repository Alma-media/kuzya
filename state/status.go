package state

import (
	"fmt"
	"strings"
)

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

// Parse converts an absolute state payload ("ON"/"OFF", case and
// whitespace tolerant) to a bool. Anything else — e.g. the transient TRIG
// payload — is not a state and is rejected.
func Parse(requested string) (bool, error) {
	switch strings.ToUpper(strings.TrimSpace(requested)) {
	case ON:
		return true, nil
	case OFF:
		return false, nil
	default:
		return false, fmt.Errorf("invalid state %q", requested)
	}
}
