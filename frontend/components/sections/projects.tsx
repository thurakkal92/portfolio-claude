import { ExternalLink, Code2, Check, TrendingUp, Zap, Star, type LucideIcon } from "lucide-react";
import { getTranslations } from "next-intl/server";
import type { Project } from "@/lib/types";

const highlightIcons: Record<string, LucideIcon> = {
  check: Check,
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
  return (
    <section
      id="projects"
      className="py-24 border-t border-border space-y-24 scroll-mt-20"
      aria-labelledby="projects-title"
    >
      <h2 id="projects-title" className="font-display text-3xl md:text-4xl tracking-tighter">
        Selected Projects
      </h2>
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
                onError={undefined}
              />
            </div>
            <div
              className={
                "md:col-span-5 space-y-4 " +
                (reversed ? "md:order-1" : "")
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
                  const Icon = highlightIcons[h.icon] ?? Check;
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
                    {t("liveDemo")} <ExternalLink className="h-3.5 w-3.5" aria-hidden />
                  </a>
                ) : null}
                {project.sourceUrl ? (
                  <a
                    href={project.sourceUrl}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="inline-flex items-center gap-1.5 font-mono text-[12px] uppercase tracking-[0.08em] underline underline-offset-4"
                  >
                    {t("source")} <Code2 className="h-3.5 w-3.5" aria-hidden />
                  </a>
                ) : null}
              </div>
            </div>
          </article>
        );
      })}
    </section>
  );
}
