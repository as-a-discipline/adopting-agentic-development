# `task service:generate`

## What it does

Delegates directly to [`api:generate`](./api-generate.md) (via
`task --taskfile ../api/Taskfile.yml --dir ../api generate`) — it does not
duplicate the generation logic.

## Why it exists

Lets you run generation from the `service` namespace without needing to
know it actually lives in `api`, while keeping exactly one implementation of
"how do we generate", per the containerized-tooling and component-boundary
ADRs.

## How to use it

```bash
task service:generate
```

## Related

* [Task Reference index](../tasks-reference.md)
* [`task api:generate`](./api-generate.md)
* [`task service:validate`](./service-validate.md)
