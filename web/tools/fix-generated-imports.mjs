#!/usr/bin/env node
// Post-processes openapi-generator-cli's "typescript-fetch" output: its relative
// import/export specifiers omit file extensions (e.g. `from '../runtime'`),
// which Node's ESM resolver requires. This script deterministically rewrites
// them to explicit `.ts` specifiers, in place, immediately after generation.
// Idempotent: safe to run multiple times.
import { readdirSync, readFileSync, statSync, writeFileSync } from "node:fs";
import { join, extname } from "node:path";

const targetDir = process.argv[2];
if (!targetDir) {
  console.error("usage: fix-generated-imports.mjs <dir>");
  process.exit(2);
}

const SPECIFIER_RE = /((?:from|import)\s+['"])(\.[^'"]+)(['"])/g;

function listTsFiles(dir) {
  const out = [];
  for (const entry of readdirSync(dir)) {
    const full = join(dir, entry);
    const st = statSync(full);
    if (st.isDirectory()) {
      out.push(...listTsFiles(full));
    } else if (extname(full) === ".ts") {
      out.push(full);
    }
  }
  return out;
}

let changed = 0;
for (const file of listTsFiles(targetDir)) {
  const original = readFileSync(file, "utf8");
  const rewritten = original.replace(SPECIFIER_RE, (match, pre, specifier, post) => {
    if (/\.[a-zA-Z]+$/.test(specifier)) {
      // already has an extension (e.g. .ts) — leave as-is.
      return match;
    }
    return `${pre}${specifier}.ts${post}`;
  });
  if (rewritten !== original) {
    writeFileSync(file, rewritten);
    changed += 1;
  }
}

console.log(`fix-generated-imports: rewrote ${changed} file(s) under ${targetDir}`);
