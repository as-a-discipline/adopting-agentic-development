import type { MonitoredService } from "../generated/models/MonitoredService.ts";
import { ServiceStatus } from "../generated/models/ServiceStatus.ts";

/**
 * Pure presentation/summary logic for the Pulse web UI, kept separate from DOM
 * rendering so it can be unit tested without a browser environment.
 */

export interface StatusPresentation {
  label: string;
  cssClass: string;
}

const STATUS_PRESENTATION: Record<string, StatusPresentation> = {
  [ServiceStatus.Healthy]: { label: "Healthy", cssClass: "status-healthy" },
  [ServiceStatus.Degraded]: { label: "Degraded", cssClass: "status-degraded" },
  [ServiceStatus.Down]: { label: "Down", cssClass: "status-down" },
  [ServiceStatus.Unknown]: { label: "Unknown", cssClass: "status-unknown" },
};

const FALLBACK_PRESENTATION: StatusPresentation = {
  label: "Unknown",
  cssClass: "status-unknown",
};

/** Maps a raw service status string to a human label and a CSS class. */
export function presentStatus(status: string): StatusPresentation {
  return STATUS_PRESENTATION[status] ?? FALLBACK_PRESENTATION;
}

export interface Summary {
  total: number;
  healthy: number;
  degraded: number;
  down: number;
  unknown: number;
}

/** Computes aggregate counts per status across a list of services. */
export function summarize(services: MonitoredService[]): Summary {
  const summary: Summary = {
    total: services.length,
    healthy: 0,
    degraded: 0,
    down: 0,
    unknown: 0,
  };

  for (const service of services) {
    switch (service.status) {
      case ServiceStatus.Healthy:
        summary.healthy += 1;
        break;
      case ServiceStatus.Degraded:
        summary.degraded += 1;
        break;
      case ServiceStatus.Down:
        summary.down += 1;
        break;
      default:
        summary.unknown += 1;
        break;
    }
  }

  return summary;
}

/** Renders a one-line human summary, e.g. "3 healthy, 1 degraded, 0 down". */
export function summaryLine(summary: Summary): string {
  if (summary.total === 0) {
    return "No monitored services.";
  }
  return `${summary.healthy} healthy, ${summary.degraded} degraded, ${summary.down} down, ${summary.unknown} unknown`;
}

/** Formats a response time in milliseconds for display, or a placeholder if absent. */
export function formatResponseTime(ms: number | null | undefined): string {
  if (ms === null || ms === undefined) {
    return "—";
  }
  return `${ms} ms`;
}

/** Formats a last-checked timestamp for display, or a placeholder if never checked. */
export function formatLastChecked(value: Date | null | undefined): string {
  if (!value) {
    return "never checked";
  }
  return value.toLocaleString();
}
