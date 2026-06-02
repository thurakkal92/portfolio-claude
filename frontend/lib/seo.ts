import type { Metadata } from "next";
import { routing } from "@/i18n/routing";
import type { ContentPayload, SocialLink } from "./types";

const SITE_URL = process.env.NEXT_PUBLIC_SITE_URL ?? "https://thurakkal.com";

export function siteUrl(): string {
  return SITE_URL.replace(/\/$/, "");
}

export function buildMetadata(content: ContentPayload, locale: string): Metadata {
  const url = `${siteUrl()}/${locale}`;
  const alternates: Record<string, string> = {};
  for (const l of routing.locales) {
    alternates[l] = `${siteUrl()}/${l}`;
  }
  return {
    metadataBase: new URL(siteUrl()),
    title: content.siteSettings.siteTitle,
    description: content.siteSettings.siteDescription,
    alternates: {
      canonical: url,
      languages: { ...alternates, "x-default": `${siteUrl()}/en` },
    },
    openGraph: {
      type: "website",
      url,
      title: content.siteSettings.siteTitle,
      description: content.siteSettings.siteDescription,
      siteName: "Nabeel Thurakkal",
      locale: locale === "de" ? "de_DE" : "en_US",
      images: content.siteSettings.ogImagePath
        ? [{ url: content.siteSettings.ogImagePath, width: 1200, height: 630 }]
        : undefined,
    },
    twitter: {
      card: "summary_large_image",
      title: content.siteSettings.siteTitle,
      description: content.siteSettings.siteDescription,
      images: content.siteSettings.ogImagePath ? [content.siteSettings.ogImagePath] : undefined,
    },
    robots: { index: true, follow: true },
  };
}

export function personJsonLd(content: ContentPayload): string {
  const sameAs = content.socialLinks
    .filter((l: SocialLink) => l.kind !== "email")
    .map((l) => l.href);
  const email = content.socialLinks.find((l) => l.kind === "email")?.href;
  const ld = {
    "@context": "https://schema.org",
    "@type": "Person",
    name: content.hero.name,
    jobTitle: "Senior Frontend Developer",
    url: siteUrl(),
    email: email?.replace(/^mailto:/, ""),
    address: {
      "@type": "PostalAddress",
      addressLocality: "Ulm",
      addressCountry: "DE",
    },
    sameAs,
  };
  return JSON.stringify(ld);
}
