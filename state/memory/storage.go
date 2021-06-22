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

func (sw *Switch) Status(deviceID string) (string, error) {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	curr := sw.states[deviceID]
	if curr {
		return state.ON, nil
	}

	return state.OFF, nil
}
