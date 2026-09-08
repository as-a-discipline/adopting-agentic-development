#!/usr/bin/env node
// Deterministic fake HTTP target for the Pulse Docker Compose integration
// environment (deploy/compose/AGENTS.md — "no external network dependencies
// for test success"). A single zero-dependency Node script, parameterized by
// environment variables, plays one of three roles depending on which
// Compose service invokes it (see docker-compose.yml):
//
//   MODE=healthy   -> GET / responds 200 immediately
//   MODE=degraded  -> GET / responds 200 after a configurable delay
//                      (DELAY_MS, default 400ms — intentionally above
//                      the service's HealthyThreshold of 300ms)
//   MODE=down      -> GET / responds 500 immediately
//
// GET /__health always responds 200 immediately, regardless of MODE — this is
// the endpoint Compose's own healthcheck uses to confirm the *container
// process* is alive, independent of the demo status it's deliberately
// simulating on `/`.
import { createServer } from "node:http";

const mode = process.env.MODE ?? "healthy";
const delayMs = Number(process.env.DELAY_MS ?? "400");
const port = Number(process.env.PORT ?? "80");

if (!["healthy", "degraded", "down"].includes(mode)) {
  console.error(`fake-target: unknown MODE '${mode}' (expected healthy|degraded|down)`);
  process.exit(2);
}

function respond(res) {
  if (mode === "down") {
    res.writeHead(500, { "content-type": "text/plain" });
    res.end("down\n");
    return;
  }
  res.writeHead(200, { "content-type": "text/plain" });
  res.end(`${mode}\n`);
}

const server = createServer((req, res) => {
  if (req.url === "/__health") {
    res.writeHead(200, { "content-type": "text/plain" });
    res.end("ok\n");
    return;
  }

  if (mode === "degraded") {
    setTimeout(() => respond(res), delayMs);
    return;
  }

  respond(res);
});

server.listen(port, () => {
  console.log(`fake-target: mode=${mode} delayMs=${delayMs} listening on :${port}`);
});
