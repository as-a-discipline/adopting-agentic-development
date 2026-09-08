# `make bootstrap`

## What it does

Installs `task` if missing (via Homebrew, then `go install`, whichever is
available), then runs [`check-prereqs`](./make-check-prereqs.md).

## Why it exists

A single command to go from a fresh clone with nothing but `git`/`docker`
installed to a machine ready to run `task validate`.

## How to use it

```bash
make bootstrap
```

## Related

* [Task Reference index](../tasks-reference.md)
* [`make check-prereqs`](./make-check-prereqs.md)
* [`task validate`](./validate.md)
