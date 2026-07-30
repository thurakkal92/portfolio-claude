"use client";

import * as React from "react";
import { Menu, X } from "lucide-react";
import { useTranslations } from "next-intl";

type NavLink = { href: string; label: string };

export function MobileMenu({ links }: { links: NavLink[] }) {
  const [open, setOpen] = React.useState(false);
  const t = useTranslations("nav");

  React.useEffect(() => {
    if (!open) return;
    const prev = document.body.style.overflow;
    document.body.style.overflow = "hidden";
    return () => {
      document.body.style.overflow = prev;
    };
  }, [open]);

  React.useEffect(() => {
    function onKey(e: KeyboardEvent) {
      if (e.key === "Escape") setOpen(false);
    }
    document.addEventListener("keydown", onKey);
    return () => document.removeEventListener("keydown", onKey);
  }, []);

  return (
    <>
      <button
        type="button"
        onClick={() => setOpen(true)}
        aria-label={t("openMenu")}
        aria-expanded={open}
        aria-controls="mobile-menu-drawer"
        className="inline-flex h-9 w-9 items-center justify-center text-muted-foreground hover:text-foreground transition-colors md:hidden"
      >
        <Menu className="h-5 w-5" aria-hidden />
      </button>

      {open ? (
        <div
          id="mobile-menu-drawer"
          role="dialog"
          aria-modal="true"
          aria-label={t("openMenu")}
          className="fixed inset-0 z-[60] bg-background md:hidden animate-fade-up"
        >
          <div className="flex justify-end px-4 h-16 items-center">
            <button
              type="button"
              onClick={() => setOpen(false)}
              aria-label={t("closeMenu")}
              className="inline-flex h-9 w-9 items-center justify-center text-muted-foreground hover:text-foreground transition-colors"
            >
              <X className="h-5 w-5" aria-hidden />
            </button>
          </div>
          <nav className="flex flex-col px-4 pt-4" aria-label="Mobile primary">
            {links.map((l) => (
              <a
                key={l.href}
                href={l.href}
                onClick={() => setOpen(false)}
                className="font-display text-xl bg-background tracking-tighter border-b border-border py-5"
              >
                {l.label}
              </a>
            ))}
          </nav>
        </div>
      ) : null}
    </>
  );
}
