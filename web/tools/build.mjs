#!/usr/bin/env node
// Zero-dependency browser build for the Pulse web app.
//
// Uses Node's built-in TypeScript type-stripping (node:module
// stripTypeScriptTypes) to compile src/**/*.ts and generated/**/*.ts into plain
// JavaScript under dist/, rewriting relative import/export specifiers from
// `.ts` to `.js` so the output is directly loadable by browsers via native
// ES modules — no bundler, no npm dependency. Also copies index.html and
// non-TS static assets (e.g. style.css) into dist/, preserving relative paths
// so index.html's `./src/...` references resolve unchanged.
import { stripTypeScriptTypes } from "node:module";
import {
  copyFileSync,
  mkdirSync,
  readdirSync,
  readFileSync,
  rmSync,
  statSync,
  writeFileSync,
} from "node:fs";
import { dirname, extname, join, relative } from "node:path";

const webDir = process.cwd();
const distDir = join(webDir, "dist");

rmSync(distDir, { recursive: true, force: true });
mkdirSync(distDir, { recursive: true });

const SPECIFIER_RE = /((?:from|import)\s+['"])(\.[^'"]+)\.ts(['"])/g;

function walk(dir) {
  const out = [];
  for (const entry of readdirSync(dir)) {
    const full = join(dir, entry);
    const st = statSync(full);
    if (st.isDirectory()) {
      out.push(...walk(full));
    } else {
      out.push(full);
    }
  }
  return out;
}

function ensureDirFor(path) {
  mkdirSync(dirname(path), { recursive: true });
}

let compiled = 0;
let copied = 0;

for (const dirName of ["src", "generated"]) {
  const sourceDir = join(webDir, dirName);
  let files;
  try {
    files = walk(sourceDir);
  } catch {
    continue; // directory doesn't exist yet — nothing to build from it.
  }

  for (const file of files) {
    const rel = relative(webDir, file);
    if (file.endsWith(".test.ts")) {
      continue; // tests are not part of the browser build.
    }
    if (extname(file) === ".ts") {
      const source = readFileSync(file, "utf8");
      const stripped = stripTypeScriptTypes(source, { mode: "transform" });
      const rewritten = stripped.replace(SPECIFIER_RE, (_m, pre, specifier, post) => `${pre}${specifier}.js${post}`);
      const outPath = join(distDir, rel.replace(/\.ts$/, ".js"));
      ensureDirFor(outPath);
      writeFileSync(outPath, rewritten);
      compiled += 1;
    } else {
      // Copy any non-.ts static asset (e.g. style.css) as-is.
      const outPath = join(distDir, rel);
      ensureDirFor(outPath);
      copyFileSync(file, outPath);
      copied += 1;
    }
  }
}

// index.html is the entry point; copy it to the dist root unchanged.
copyFileSync(join(webDir, "index.html"), join(distDir, "index.html"));

console.log(`build: compiled ${compiled} TypeScript file(s), copied ${copied} static asset(s), plus index.html -> dist/`);
