import test from "node:test";
import assert from "node:assert/strict";
import type { MonitoredService } from "../generated/models/MonitoredService.ts";
import { ServiceStatus } from "../generated/models/ServiceStatus.ts";
import {
  formatLastChecked,
  formatResponseTime,
  presentStatus,
  summarize,
  summaryLine,
} from "./pulse.ts";

function service(overrides: Partial<MonitoredService> = {}): MonitoredService {
  return {
    id: "svc-1",
    name: "Example",
    url: "https://example.invalid",
    status: ServiceStatus.Unknown,
    lastChecked: null,
    responseTimeMs: null,
    message: null,
    ...overrides,
  };
}

test("presentStatus maps known statuses", () => {
  assert.deepEqual(presentStatus(ServiceStatus.Healthy), {
    label: "Healthy",
    cssClass: "status-healthy",
  });
  assert.deepEqual(presentStatus(ServiceStatus.Degraded), {
    label: "Degraded",
    cssClass: "status-degraded",
  });
  assert.deepEqual(presentStatus(ServiceStatus.Down), {
    label: "Down",
    cssClass: "status-down",
  });
});

test("presentStatus falls back to unknown for an unrecognized value", () => {
  assert.deepEqual(presentStatus("bogus"), {
    label: "Unknown",
    cssClass: "status-unknown",
  });
});

test("summarize counts each status independently", () => {
  const services = [
    service({ id: "a", status: ServiceStatus.Healthy }),
    service({ id: "b", status: ServiceStatus.Healthy }),
    service({ id: "c", status: ServiceStatus.Degraded }),
    service({ id: "d", status: ServiceStatus.Down }),
    service({ id: "e", status: ServiceStatus.Unknown }),
  ];

  assert.deepEqual(summarize(services), {
    total: 5,
    healthy: 2,
    degraded: 1,
    down: 1,
    unknown: 1,
  });
});

test("summarize on an empty list returns all zero counts", () => {
  assert.deepEqual(summarize([]), {
    total: 0,
    healthy: 0,
    degraded: 0,
    down: 0,
    unknown: 0,
  });
});

test("summaryLine reports counts for a non-empty list", () => {
  const summary = summarize([
    service({ status: ServiceStatus.Healthy }),
    service({ status: ServiceStatus.Down }),
  ]);
  assert.equal(summaryLine(summary), "1 healthy, 0 degraded, 1 down, 0 unknown");
});

test("summaryLine reports a placeholder message for an empty list", () => {
  assert.equal(summaryLine(summarize([])), "No monitored services.");
});

test("formatResponseTime formats a present value", () => {
  assert.equal(formatResponseTime(42), "42 ms");
});

test("formatResponseTime formats null/undefined as a placeholder", () => {
  assert.equal(formatResponseTime(null), "—");
  assert.equal(formatResponseTime(undefined), "—");
});

test("formatLastChecked formats a present date", () => {
  const date = new Date("2024-01-01T00:00:00.000Z");
  assert.equal(formatLastChecked(date), date.toLocaleString());
});

test("formatLastChecked formats null/undefined as never-checked", () => {
  assert.equal(formatLastChecked(null), "never checked");
  assert.equal(formatLastChecked(undefined), "never checked");
});
