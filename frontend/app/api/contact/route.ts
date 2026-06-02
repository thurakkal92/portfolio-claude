import { NextResponse } from "next/server";

const API_BASE = process.env.API_BASE_URL ?? "http://localhost:8080";

export async function POST(req: Request) {
  const body = await req.text();
  // Forward client IP & UA so the Go service can rate-limit / audit accurately.
  const forwardedFor =
    req.headers.get("x-forwarded-for") ??
    req.headers.get("x-real-ip") ??
    "";
  const userAgent = req.headers.get("user-agent") ?? "";

  const upstream = await fetch(`${API_BASE}/api/contact`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "X-Forwarded-For": forwardedFor,
      "User-Agent": userAgent,
    },
    body,
  });

  const text = await upstream.text();
  return new NextResponse(text, {
    status: upstream.status,
    headers: { "Content-Type": upstream.headers.get("Content-Type") ?? "application/json" },
  });
}
