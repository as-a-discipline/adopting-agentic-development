import { Configuration, ServicesApi } from "../generated/index.ts";
import type { MonitoredService } from "../generated/models/MonitoredService.ts";
import {
  formatLastChecked,
  formatResponseTime,
  presentStatus,
  summarize,
  summaryLine,
} from "./pulse.ts";

// Runtime-configurable API base URL: index.html may set `window.PULSE_API_BASE_URL`
// before loading this script (the Compose environment does this to point at the
// pulse-service container). Defaults to the same-origin "/api/v1" path.
declare global {
  interface Window {
    PULSE_API_BASE_URL?: string;
  }
}
const apiBaseUrl = window.PULSE_API_BASE_URL ?? "/api/v1";

const api = new ServicesApi(new Configuration({ basePath: apiBaseUrl }));

const app = document.querySelector<HTMLDivElement>("#app")!;

function render(services: MonitoredService[], error?: string): void {
  const summary = summarize(services);

  app.innerHTML = `
    <header>
      <h1>Pulse</h1>
      <p class="summary">${summaryLine(summary)}</p>
      ${error ? `<p class="error" role="alert">${escapeHtml(error)}</p>` : ""}
    </header>
    <ul class="service-list">
      ${services.map(renderCard).join("")}
    </ul>
  `;

  for (const service of services) {
    const button = app.querySelector<HTMLButtonElement>(
      `[data-check-id="${service.id}"]`
    );
    button?.addEventListener("click", () => onCheck(service.id));
  }
}

function renderCard(service: MonitoredService): string {
  const { label, cssClass } = presentStatus(service.status);
  return `
    <li class="service-card ${cssClass}">
      <h2>${escapeHtml(service.name)}</h2>
      <p class="status" aria-label="status">${label}</p>
      <dl>
        <dt>URL</dt>
        <dd>${escapeHtml(service.url)}</dd>
        <dt>Response time</dt>
        <dd>${formatResponseTime(service.responseTimeMs)}</dd>
        <dt>Last checked</dt>
        <dd>${formatLastChecked(service.lastChecked)}</dd>
      </dl>
      <button type="button" data-check-id="${service.id}">Check now</button>
    </li>
  `;
}

function escapeHtml(value: string): string {
  const div = document.createElement("div");
  div.textContent = value;
  return div.innerHTML;
}

async function loadServices(): Promise<void> {
  try {
    const services = await api.listServices();
    render(services);
  } catch (err) {
    render([], `Failed to load services: ${(err as Error).message}`);
  }
}

async function onCheck(id: string): Promise<void> {
  try {
    await api.checkService({ id });
  } catch {
    // The list reload below will surface any resulting error state; a failed
    // check request itself is not fatal to the UI.
  }
  await loadServices();
}

loadServices();
