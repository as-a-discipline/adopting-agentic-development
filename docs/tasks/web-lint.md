# `task web:lint`

## What it does

Runs `tools/check.mjs` in a pinned `node:22-alpine` container — a syntax +
import-resolution check, not full TypeScript type checking, since no `tsc`
binary is installed (see `web/AGENTS.md` for why).

## Why it exists

Catches syntax errors and broken import paths in web source without adding
an npm dependency — this repo's web toolchain is deliberately
dependency-free (see `web/AGENTS.md`).

## How to use it

```bash
task web:lint
```

## Related

* [Task Reference index](../tasks-reference.md)
* [`task web:test`](./web-test.md)
* [`task web:validate`](./web-validate.md)
