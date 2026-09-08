# `task structure:validate`

## What it does

Runs [`scripts/validate-structure.sh`](../../scripts/validate-structure.sh),
which checks that:

* required root files exist,
* every component has an `AGENTS.md`,
* every component has a `Taskfile.yml`,
* all 6 skills exist under `.agents/skills/` with valid `SKILL.md`
  frontmatter (`name`/`description`),
* the required root ADR directory and each component's `adr/` directory
  exist,
* the root verification docs (`verification/skills.md`,
  `verification/architecture.md`) exist.

## Why it exists

A deterministic check that the repository's own documented shape
(instructions, skills, task interfaces, ADRs) hasn't silently rotted. It is
the first, fastest thing `task validate` checks, and it's also one of the
fast checks run by the local pre-commit hook (see
[`docs/api-breaking-changes.md`](../api-breaking-changes.md) for how the
hook is wired).

## How to use it

```bash
task structure:validate
```

## Related

* [Task Reference index](../tasks-reference.md)
* [`task validate`](./validate.md)
