import { Code, Layers, Wrench, Terminal, type LucideIcon } from "lucide-react";
import type { SkillGroup } from "@/lib/types";

const icons: Record<string, LucideIcon> = {
  code: Code,
  layers: Layers,
  wrench: Wrench,
  terminal: Terminal,
};

export function SkillsSection({ skills }: { skills: SkillGroup[] }) {
  return (
    <section
      className="py-16 md:py-24 border-t border-border"
      aria-labelledby="skills-title"
    >
      <h2 id="skills-title" className="font-display text-2xl md:text-4xl tracking-tighter mb-8 md:mb-10">
        Core Competencies
      </h2>
      <div className="grid grid-cols-2 md:grid-cols-2 gap-3 md:gap-6">
        {skills.map((group) => {
          const Icon = icons[group.icon] ?? Code;
          return (
            <div
              key={group.slug}
              className="border border-border bg-card p-4 md:p-6 flex flex-col gap-3 md:gap-4"
            >
              <div className="flex flex-col md:flex-row md:items-center gap-2 md:gap-2">
                <Icon className="h-5 w-5" aria-hidden />
                <h3 className="font-display text-base md:text-lg font-semibold">
                  {group.title}
                </h3>
              </div>
              <ul className="hidden md:flex flex-wrap gap-2" role="list">
                {group.items.map((item) => (
                  <li
                    key={item}
                    className="bg-muted px-2.5 py-1 font-mono text-[11px] tracking-[0.04em] text-foreground/90"
                  >
                    {item}
                  </li>
                ))}
              </ul>
              <p className="md:hidden text-xs text-muted-foreground leading-relaxed">
                {group.items.slice(0, 3).join(" · ")}
                {group.items.length > 3 ? " · …" : ""}
              </p>
            </div>
          );
        })}
      </div>
    </section>
  );
}
