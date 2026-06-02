import "server-only";
import type { ContentPayload } from "./types";

const API_BASE = process.env.API_BASE_URL ?? "http://localhost:8080";

export async function fetchContent(locale: string): Promise<ContentPayload> {
  const url = `${API_BASE}/api/content?locale=${encodeURIComponent(locale)}`;
  const res = await fetch(url, {
    next: { revalidate: 60 },
  });
  if (!res.ok) {
    const body = await res.text().catch(() => "");
    throw new Error(`fetchContent ${locale}: ${res.status} ${body}`);
  }
  return (await res.json()) as ContentPayload;
}

export type ContactPayload = {
  name: string;
  email: string;
  message: string;
  locale: string;
  website?: string;
};

export type ContactResult =
  | { ok: true }
  | { ok: false; status: number; error?: string; fields?: Record<string, string> };

// Server-side proxy; called from app/api/contact/route.ts.
export async function submitContact(payload: ContactPayload): Promise<ContactResult> {
  const res = await fetch(`${API_BASE}/api/contact`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });
  if (res.ok) return { ok: true };
  let body: { error?: string; fields?: Record<string, string> } = {};
  try {
    body = await res.json();
  } catch {
    /* ignore */
  }
  return { ok: false, status: res.status, error: body.error, fields: body.fields };
}
