package memory

import (
	"sync"

	"github.com/Alma-media/kuzya/state"
)

type Switch struct {
	mu     sync.Mutex
	states map[string]bool
}

func NewSwitch() *Switch {
	return &Switch{
		states: make(map[string]bool),
	}
}

func (sw *Switch) Switch(deviceID string) (string, error) {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	curr := sw.states[deviceID]
	if !curr {
		sw.states[deviceID] = true

		return state.ON, nil
	}

	sw.states[deviceID] = false

	return state.OFF, nil
}

// SetState applies an absolute state ("ON"/"OFF") instead of toggling.
// Scenes and schedules need it: toggling an unknown-phase channel is how
// "evening scene" ends up turning the lights off.
func (sw *Switch) SetState(deviceID, requested string) (string, error) {
	target, err := state.Parse(requested)
	if err != nil {
		return "", err
	}

	sw.mu.Lock()
	defer sw.mu.Unlock()

	sw.states[deviceID] = target

	return state.Status(target), nil
}

func (sw *Switch) Status(deviceID string) (string, error) {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	curr := sw.states[deviceID]
	if curr {
		return state.ON, nil
	}

	return state.OFF, nil
}
