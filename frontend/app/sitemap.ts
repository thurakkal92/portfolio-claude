import type { MetadataRoute } from "next";
import { routing } from "@/i18n/routing";
import { siteUrl } from "@/lib/seo";

export default function sitemap(): MetadataRoute.Sitemap {
  const base = siteUrl();
  const now = new Date();
  const entries: MetadataRoute.Sitemap = [];
  for (const locale of routing.locales) {
    const url = `${base}/${locale}`;
    const alternates: Record<string, string> = {};
    for (const l of routing.locales) alternates[l] = `${base}/${l}`;
    entries.push({
      url,
      lastModified: now,
      changeFrequency: "monthly",
      priority: locale === routing.defaultLocale ? 1.0 : 0.8,
      alternates: { languages: alternates },
    });
  }
  return entries;
}
