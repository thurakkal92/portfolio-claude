import { getTranslations } from "next-intl/server";
import type { Experience } from "@/lib/types";
import { formatDateRange } from "@/lib/utils";

export async function ExperienceSection({
  experiences,
  locale,
}: {
  experiences: Experience[];
  locale: string;
}) {
  const t = await getTranslations({ locale, namespace: "experience" });
  const present = t("present");
  return (
    <section
      id="experience"
      className="py-24 border-t border-border scroll-mt-20"
      aria-labelledby="experience-title"
    >
      <h2 id="experience-title" className="font-display text-3xl md:text-4xl tracking-tighter mb-12">
        Professional Experience
      </h2>
      <ol
        className="relative space-y-12 before:content-[''] before:absolute before:left-2 md:before:left-1/2 before:top-2 before:bottom-2 before:w-px before:bg-border"
        aria-label="Career timeline"
      >
        {experiences.map((xp, i) => {
          const right = i % 2 === 1;
          const range = formatDateRange(xp.startDate, xp.endDate, locale, present);
          const meta = [range, xp.location].filter(Boolean).join(" · ");
          return (
            <li key={i} className="relative pl-8 md:pl-0">
              <span
                aria-hidden
                className={
                  "absolute top-2 h-3 w-3 rounded-full border-4 border-background z-10 " +
                  (i === 0 ? "bg-foreground" : "bg-border") +
                  " left-0.5 md:left-1/2 md:-translate-x-1/2"
                }
              />
              <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                <div className={right ? "md:col-start-2 md:pl-10" : "md:pr-10 md:text-right"}>
                  <h3 className="font-display text-xl">{xp.company}</h3>
                  <p className="label-mono-md text-muted-foreground mt-1">{xp.role}</p>
                  <p className="font-mono text-[11px] uppercase tracking-[0.08em] text-muted-foreground mt-1">
                    {meta}
                  </p>
                </div>
                <div
                  className={
                    right ? "md:col-start-1 md:row-start-1 md:pr-10 md:text-right" : ""
                  }
                >
                  <div className="border border-border bg-card p-5">
                    <p className="text-sm text-muted-foreground leading-relaxed text-left">
                      {xp.summary}
                    </p>
                  </div>
                </div>
              </div>
            </li>
          );
        })}
      </ol>
    </section>
  );
}
