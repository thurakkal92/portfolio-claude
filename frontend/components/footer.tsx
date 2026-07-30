import { getTranslations } from "next-intl/server";
import { Pi } from "lucide-react";

export async function Footer({ locale }: { locale: string }) {
  const t = await getTranslations({ locale, namespace: "footer" });
  const nav = await getTranslations({ locale, namespace: "nav" });
  const year = new Date().getFullYear();
  return (
    <footer className="border-t border-border bg-background mt-24">
      <div className="container mx-auto py-16 grid grid-cols-1 md:grid-cols-3 gap-8">
        <div className="space-y-4">
          <div className="text-foreground" aria-label="Nabeel Thurakkal">
            <Pi className="h-6 w-6" strokeWidth={2.25} aria-hidden />
          </div>
          <p className="text-sm text-muted-foreground max-w-xs">{t("tagline")}</p>
        </div>
        <div className="flex flex-col gap-2">
          <h4 className="label-mono mb-1">{t("quickLinks")}</h4>
          <a href="#about" className="text-sm text-muted-foreground hover:text-foreground">{nav("about")}</a>
          <a href="#projects" className="text-sm text-muted-foreground hover:text-foreground">{nav("projects")}</a>
          <a href="#experience" className="text-sm text-muted-foreground hover:text-foreground">{nav("experience")}</a>
        </div>
        <div className="flex flex-col gap-2">
          <h4 className="label-mono mb-1">{t("legal")}</h4>
          <a href={`/${locale}/impressum`} className="text-sm text-muted-foreground hover:text-foreground">
            {t("impressum")}
          </a>
          <a href={`/${locale}/datenschutz`} className="text-sm text-muted-foreground hover:text-foreground">
            {t("datenschutz")}
          </a>
        </div>
      </div>
      <div className="container mx-auto pb-8 pt-4 border-t border-border/50 flex items-center justify-between">
        <p className="text-xs text-muted-foreground font-mono">
          {t("copyright", { year })}
        </p>
      </div>
    </footer>
  );
}
