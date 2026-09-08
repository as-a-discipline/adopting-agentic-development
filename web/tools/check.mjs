#!/usr/bin/env node
// Dependency-free static check for the web app (stands in for a linter/type
// checker without requiring an npm-installed toolchain — see AGENTS.md).
//
// For every .ts file under src/ and generated/:
//   1. Confirms it parses/strips cleanly via Node's built-in TypeScript type
//      stripping (catches syntax errors).
//   2. Confirms every relative import/export specifier resolves to a real
//      file on disk (catches broken imports / missing extensions).
//
// This does not perform full type checking (no TypeScript compiler is
// installed by design); it is a syntax + import-resolution safety net.
import { stripTypeScriptTypes } from "node:module";
import { existsSync, readdirSync, readFileSync, statSync } from "node:fs";
import { dirname, extname, join, relative, resolve } from "node:path";

const webDir = process.cwd();
const SPECIFIER_RE = /(?:from|import)\s+['"](\.[^'"]+)['"]/g;

function walk(dir) {
  const out = [];
  for (const entry of readdirSync(dir)) {
    const full = join(dir, entry);
    const st = statSync(full);
    if (st.isDirectory()) {
      out.push(...walk(full));
    } else if (extname(full) === ".ts") {
      out.push(full);
    }
  }
  return out;
}

const errors = [];
let checked = 0;

for (const dirName of ["src", "generated"]) {
  const sourceDir = join(webDir, dirName);
  if (!existsSync(sourceDir)) continue;

  for (const file of walk(sourceDir)) {
    checked += 1;
    const rel = relative(webDir, file);
    const source = readFileSync(file, "utf8");

    try {
      stripTypeScriptTypes(source, { mode: "transform" });
    } catch (err) {
      errors.push(`${rel}: syntax error — ${err.message}`);
      continue;
    }

    for (const match of source.matchAll(SPECIFIER_RE)) {
      const specifier = match[1];
      const resolved = resolve(dirname(file), specifier);
      if (!existsSync(resolved)) {
        errors.push(`${rel}: import '${specifier}' does not resolve to an existing file`);
      }
    }
  }
}

if (errors.length > 0) {
  console.error(`check: ${errors.length} problem(s) found across ${checked} file(s):`);
  for (const error of errors) {
    console.error(`  - ${error}`);
  }
  process.exit(1);
}

console.log(`check: OK — ${checked} file(s) passed syntax and import-resolution checks.`);
