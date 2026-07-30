import {
  ExternalLink,
  Code2,
  CheckCircle2,
  TrendingUp,
  Zap,
  Star,
  ArrowUpRight,
  type LucideIcon,
} from "lucide-react";
import { getTranslations } from "next-intl/server";
import type { Project } from "@/lib/types";

const highlightIcons: Record<string, LucideIcon> = {
  check: CheckCircle2,
  "trending-up": TrendingUp,
  zap: Zap,
  star: Star,
};

export async function ProjectsSection({
  projects,
  locale,
}: {
  projects: Project[];
  locale: string;
}) {
  const t = await getTranslations({ locale, namespace: "projects" });
  const total = String(projects.length).padStart(2, "0");
  return (
    <section
      id="projects"
      className="py-16 md:py-24 border-t border-border space-y-10 md:space-y-24 scroll-mt-20"
      aria-labelledby="projects-title"
    >
      <div className="flex items-end justify-between gap-4">
        <h2
          id="projects-title"
          className="font-display text-2xl md:text-4xl tracking-tighter"
        >
          Selected Projects
        </h2>
        <span
          className="md:hidden font-mono text-[11px] uppercase tracking-[0.12em] text-muted-foreground"
          aria-hidden
        >
          01 — {total}
        </span>
      </div>

      {/* Desktop: alternating grid, hidden on mobile */}
      <div className="hidden md:block space-y-24">
        {projects.map((project, idx) => {
          const reversed = idx % 2 === 1;
          return (
            <article
              key={project.slug}
              className="grid grid-cols-1 md:grid-cols-12 gap-8 items-center"
            >
              <div
                className={
                  "md:col-span-7 overflow-hidden border border-border bg-muted " +
                  (reversed ? "md:order-2" : "")
                }
              >
                {/* eslint-disable-next-line @next/next/no-img-element */}
                <img
                  src={project.imagePath}
                  alt={`${project.title} — preview`}
                  loading="lazy"
                  className="w-full aspect-video object-cover grayscale hover:grayscale-0 transition-all duration-500 motion-reduce:transition-none"
                />
              </div>
              <div
                className={
                  "md:col-span-5 space-y-4 " + (reversed ? "md:order-1" : "")
                }
              >
                <p className="label-mono">{project.company}</p>
                <h3 className="font-display text-2xl md:text-3xl tracking-tighter">
                  {project.title}
                </h3>
                <p className="text-base text-muted-foreground leading-relaxed">
                  {project.description}
                </p>
                <ul className="flex flex-wrap gap-2" role="list">
                  {project.tags.map((tag) => (
                    <li
                      key={tag}
                      className="bg-accent px-2 py-1 font-mono text-[11px]"
                    >
                      {tag}
                    </li>
                  ))}
                </ul>
                <ul className="space-y-1.5 text-sm">
                  {project.highlights.map((h, i) => {
                    const Icon = highlightIcons[h.icon] ?? CheckCircle2;
                    return (
                      <li key={i} className="flex items-start gap-2">
                        <Icon className="h-4 w-4 mt-0.5 shrink-0" aria-hidden />
                        <span>{h.text}</span>
                      </li>
                    );
                  })}
                </ul>
                <div className="flex gap-4 pt-2">
                  {project.liveUrl ? (
                    <a
                      href={project.liveUrl}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="inline-flex items-center gap-1.5 font-mono text-[12px] uppercase tracking-[0.08em] underline underline-offset-4"
                    >
                      {t("liveDemo")}{" "}
                      <ExternalLink className="h-3.5 w-3.5" aria-hidden />
                    </a>
                  ) : null}
                  {project.sourceUrl ? (
                    <a
                      href={project.sourceUrl}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="inline-flex items-center gap-1.5 font-mono text-[12px] uppercase tracking-[0.08em] underline underline-offset-4"
                    >
                      {t("source")}{" "}
                      <Code2 className="h-3.5 w-3.5" aria-hidden />
                    </a>
                  ) : null}
                </div>
              </div>
            </article>
          );
        })}
      </div>

      {/* Mobile: simplified vertical stack with alternating background */}
      <div className="md:hidden -mx-4">
        {projects.map((project, idx) => {
          const primaryHref = project.liveUrl || project.sourceUrl;
          const external = !!primaryHref;
          return (
            <article
              key={project.slug}
              className={
                "flex flex-col gap-4 px-4 py-10 border-b border-border last:border-b-0 " +
                (idx % 2 === 1 ? "bg-muted/50" : "")
              }
            >
              <div className="aspect-[16/9] w-full overflow-hidden border border-border bg-muted">
                {/* eslint-disable-next-line @next/next/no-img-element */}
                <img
                  src={project.imagePath}
                  alt={`${project.title} — preview`}
                  loading="lazy"
                  className="w-full h-full object-cover"
                />
              </div>
              <div className="flex flex-col gap-3">
                {project.tags.length > 0 ? (
                  <ul className="flex flex-wrap gap-2" role="list">
                    {project.tags.slice(0, 2).map((tag) => (
                      <li
                        key={tag}
                        className="bg-accent px-2 py-1 font-mono text-[11px]"
                      >
                        {tag}
                      </li>
                    ))}
                  </ul>
                ) : null}
                <h3 className="font-display text-2xl tracking-tighter">
                  {project.title}
                </h3>
                <p className="text-sm text-muted-foreground leading-relaxed">
                  {project.description}
                </p>
                {primaryHref ? (
                  <a
                    href={primaryHref}
                    target={external ? "_blank" : undefined}
                    rel={external ? "noopener noreferrer" : undefined}
                    className="mt-2 inline-flex items-center gap-1.5 font-mono text-[12px] uppercase tracking-[0.12em] text-foreground"
                  >
                    {t("viewCaseStudy")}{" "}
                    <ArrowUpRight className="h-3.5 w-3.5" aria-hidden />
                  </a>
                ) : null}
              </div>
            </article>
          );
        })}
      </div>
    </section>
  );
}
