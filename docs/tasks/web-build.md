# `task web:build`

## What it does

Runs `tools/build.mjs` (Node's built-in `stripTypeScriptTypes`) in the
pinned container to compile `.ts` sources to static, browser-loadable
JavaScript.

## Why it exists

Produces the actual deployable web asset — the same build `web/Dockerfile`
packages into the nginx image used by Compose/Helm.

## How to use it

```bash
task web:build
```

## Related

* [Task Reference index](../tasks-reference.md)
* [`task web:test`](./web-test.md)
* [`task web:validate`](./web-validate.md)
