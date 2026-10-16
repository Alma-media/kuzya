# Status manager

### Endpoints
| input | payload | effect | output (`/trig.out.…/<ID>`) |
|---|---|---|---|
| `/trig.in.switch/<ID>` | any | toggle the channel | `switch` — new state `ON`/`OFF` |
| `/trig.in.setswitch/<ID>` | `ON` / `OFF` | set the channel absolutely | `switch` — new state `ON`/`OFF` |
| `/trig.in.status/<ID>` | any | query the channel | `status` — current state `ON`/`OFF` |

`setswitch` exists for controllers that command absolute states (scenes,
schedules): toggling would make them depend on a phase they cannot see.
Its result goes to `/trig.out.switch/<ID>` — exactly like a toggle — so no
extra wiring is needed for relays, LEDs or mirrors.

### Message forwarding example
```
[
    {
        "input": "/switch/<UUID>",
        "output": "/trig.in.switch/<ID>",
        "retain": false
    },
    {
        "input": "/trig.out.switch/<ID>",
        "output": "/led/<UUID>",
        "retain": true
    },
    {
        "input": "/trig.out.switch/<ID>",
        "output": "/relay/<UUID>",
        "retain": true
    },
    . . .
]