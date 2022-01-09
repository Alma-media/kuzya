# Status manager

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