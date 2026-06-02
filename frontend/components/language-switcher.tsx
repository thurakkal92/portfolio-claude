"use client";

import * as React from "react";
import { Languages } from "lucide-react";
import { useTranslations } from "next-intl";
import { usePathname, useRouter } from "@/i18n/routing";
import { routing } from "@/i18n/routing";

export function LanguageSwitcher({ currentLocale }: { currentLocale: string }) {
  const t = useTranslations("nav");
  const router = useRouter();
  const pathname = usePathname();
  const [open, setOpen] = React.useState(false);
  const containerRef = React.useRef<HTMLDivElement>(null);

  React.useEffect(() => {
    function onClick(e: MouseEvent) {
      if (!containerRef.current?.contains(e.target as Node)) setOpen(false);
    }
    document.addEventListener("mousedown", onClick);
    return () => document.removeEventListener("mousedown", onClick);
  }, []);

  return (
    <div ref={containerRef} className="relative">
      <button
        type="button"
        onClick={() => setOpen((v) => !v)}
        aria-haspopup="listbox"
        aria-expanded={open}
        aria-label={t("switchLanguage")}
        className="inline-flex h-9 items-center gap-1.5 px-2 text-muted-foreground hover:text-foreground transition-colors"
      >
        <Languages className="h-4 w-4" />
        <span className="font-mono text-[11px] uppercase tracking-[0.12em]">
          {currentLocale}
        </span>
      </button>
      {open && (
        <ul
          role="listbox"
          className="absolute right-0 top-full mt-1 min-w-[120px] border border-border bg-card shadow-md z-50"
        >
          {routing.locales.map((l) => (
            <li key={l}>
              <button
                type="button"
                role="option"
                aria-selected={l === currentLocale}
                onClick={() => {
                  setOpen(false);
                  router.replace(pathname, { locale: l });
                }}
                className="block w-full px-3 py-2 text-left font-mono text-[12px] uppercase tracking-[0.08em] hover:bg-muted"
              >
                {l === "en" ? "English" : "Deutsch"}
              </button>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
