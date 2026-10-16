package memory

import (
	"testing"

	"github.com/Alma-media/kuzya/state"
)

func TestSetState(t *testing.T) {
	sw := NewSwitch()

	if got, err := sw.SetState("lamp", "ON"); err != nil || got != state.ON {
		t.Fatalf(`SetState("ON") = %q, %v; want %q, nil`, got, err, state.ON)
	}

	if got, err := sw.SetState("lamp", " off "); err != nil || got != state.OFF {
		t.Fatalf(`SetState(" off ") = %q, %v; want %q, nil`, got, err, state.OFF)
	}

	if got, _ := sw.Status("lamp"); got != state.OFF {
		t.Errorf("Status = %q, want %q", got, state.OFF)
	}

	if _, err := sw.SetState("lamp", "TRIG"); err == nil {
		t.Error("TRIG payload must be rejected")
	}

	if got, _ := sw.Status("lamp"); got != state.OFF {
		t.Errorf("state after rejected payload = %q, want %q", got, state.OFF)
	}

	// Toggle continues from the set value — one shared phase for buttons
	// and absolute controllers.
	if got, err := sw.Switch("lamp"); err != nil || got != state.ON {
		t.Fatalf(`Switch after set OFF = %q, %v; want %q, nil`, got, err, state.ON)
	}
}
