import type { Metadata } from "next";
import { setRequestLocale } from "next-intl/server";
import { fetchContent } from "@/lib/api";
import { buildMetadata, personJsonLd } from "@/lib/seo";
import { HeroSection } from "@/components/sections/hero";
import { AboutSection } from "@/components/sections/about";
import { SkillsSection } from "@/components/sections/skills";
import { ProjectsSection } from "@/components/sections/projects";
import { ExperienceSection } from "@/components/sections/experience";
import { ContactSection } from "@/components/sections/contact";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  const content = await fetchContent(locale);
  return buildMetadata(content, locale);
}

export default async function HomePage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  setRequestLocale(locale);
  const content = await fetchContent(locale);
  return (
    <>
      <script
        type="application/ld+json"
        // eslint-disable-next-line react/no-danger
        dangerouslySetInnerHTML={{ __html: personJsonLd(content) }}
      />
      <HeroSection hero={content.hero} site={content.siteSettings} />
      <AboutSection about={content.about} />
      <SkillsSection skills={content.skills} />
      <ProjectsSection projects={content.projects} locale={locale} />
      <ExperienceSection experiences={content.experiences} locale={locale} />
      <ContactSection
        copy={content.contact}
        socials={content.socialLinks}
        locale={locale}
      />
    </>
  );
}
