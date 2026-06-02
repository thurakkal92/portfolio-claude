import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function formatDateRange(
  startDate: string,
  endDate: string | undefined,
  locale: string,
  presentLabel: string,
): string {
  const start = formatDate(startDate, locale);
  const end = endDate ? formatDate(endDate, locale) : presentLabel;
  return `${start} – ${end}`;
}

function formatDate(iso: string, locale: string): string {
  const [y, m] = iso.split("-").map(Number);
  if (!y || !m) return iso;
  return new Intl.DateTimeFormat(locale === "de" ? "de-DE" : "en-US", {
    month: "short",
    year: "numeric",
  }).format(new Date(y, m - 1, 1));
}
