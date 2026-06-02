import { Button } from "@/components/ui/button";
import type { Hero, SiteSettings } from "@/lib/types";

export function HeroSection({ hero, site }: { hero: Hero; site: SiteSettings }) {
  return (
    <section
      className="relative py-24 md:py-32 overflow-hidden"
      aria-labelledby="hero-title"
    >
      <div
        aria-hidden
        className="absolute top-0 right-0 -z-10 w-1/2 h-full opacity-10 bg-[radial-gradient(circle_at_center,hsl(var(--foreground))_0%,transparent_60%)]"
      />
      <div className="space-y-6 animate-fade-up">
        <p className="label-mono-md text-muted-foreground">{hero.eyebrow}</p>
        <h1
          id="hero-title"
          className="font-display text-5xl md:text-6xl font-bold tracking-tightest leading-[1.05] max-w-3xl"
        >
          {hero.name}
        </h1>
        <p className="text-lg text-muted-foreground max-w-xl leading-relaxed">
          {hero.subtitle}
        </p>
        <div className="flex flex-wrap gap-3 pt-2">
          <a href="#projects">
            <Button variant="primary" size="lg">
              {hero.ctaPrimaryLabel}
            </Button>
          </a>
          {site.cvPath ? (
            <a href={site.cvPath} download>
              <Button variant="outline" size="lg">
                {hero.ctaSecondaryLabel}
              </Button>
            </a>
          ) : null}
        </div>
      </div>
    </section>
  );
}
