import type { About } from "@/lib/types";

export function AboutSection({ about }: { about: About }) {
  const paragraphs = about.bodyMd.split(/\n\n+/);
  return (
    <section
      id="about"
      className="py-24 border-t border-border scroll-mt-20"
      aria-labelledby="about-title"
    >
      <div className="grid grid-cols-1 md:grid-cols-12 gap-12">
        <div className="md:col-span-7 space-y-4">
          <h2 id="about-title" className="font-display text-3xl md:text-4xl tracking-tighter">
            {about.heading}
          </h2>
          <div className="space-y-4 text-base text-muted-foreground leading-relaxed">
            {paragraphs.map((p, i) => (
              <p key={i}>{p}</p>
            ))}
          </div>
        </div>
        <div className="md:col-span-5">
          <div className="border border-border bg-card p-6 space-y-4">
            <h3 className="label-mono">Quick Facts</h3>
            <ul className="space-y-3">
              {about.quickFacts.map((fact, i) => (
                <li
                  key={i}
                  className={
                    "flex items-center justify-between border-b border-border/40 pb-2 last:border-0 last:pb-0"
                  }
                >
                  <span className="font-mono text-xs uppercase tracking-[0.08em] text-muted-foreground">
                    {fact.label}
                  </span>
                  {fact.kind === "status" ? (
                    <span className="inline-flex items-center gap-2 text-sm font-medium">
                      <span
                        className="h-2 w-2 rounded-full bg-emerald-500"
                        aria-hidden
                      />
                      {fact.value}
                    </span>
                  ) : (
                    <span className="text-sm font-medium text-right">
                      {fact.value}
                    </span>
                  )}
                </li>
              ))}
            </ul>
          </div>
        </div>
      </div>
    </section>
  );
}
