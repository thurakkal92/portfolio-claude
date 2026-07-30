import { Mail, Linkedin, Github } from "lucide-react";
import type { ContactCopy, SocialLink } from "@/lib/types";
import { ContactForm } from "@/components/contact-form";

export function ContactSection({
  copy,
  socials,
  locale,
}: {
  copy: ContactCopy;
  socials: SocialLink[];
  locale: string;
}) {
  return (
    <section
      id="contact"
      className="pt-16 md:pt-24 pb-0 md:pb-24 border-t border-border scroll-mt-20"
      aria-labelledby="contact-title"
    >
      {/* Mobile: elevated panel with rounded top; Desktop: transparent, centered */}
      <div className="-mx-4 md:mx-0 bg-muted md:bg-transparent rounded-t-3xl md:rounded-none px-4 md:px-0 pt-10 md:pt-0 pb-16 md:pb-0">
        <div className="max-w-xl mx-auto space-y-8 md:space-y-10 text-left md:text-center">
          <div className="space-y-2">
            <h2 id="contact-title" className="font-display text-2xl md:text-4xl tracking-tighter">
              {copy.heading}
            </h2>
            <p className="text-base text-muted-foreground">{copy.subheading}</p>
          </div>
          <ContactForm copy={copy} locale={locale} />
          <ul
            className="flex flex-wrap justify-start md:justify-center gap-6 md:gap-8 pt-4 md:pt-6"
            role="list"
          >
            {socials.map((s) => {
              const Icon =
                s.kind === "email" ? Mail : s.kind === "linkedin" ? Linkedin : Github;
              return (
                <li key={s.kind}>
                  <a
                    href={s.href}
                    target={s.kind === "email" ? undefined : "_blank"}
                    rel={s.kind === "email" ? undefined : "noopener noreferrer"}
                    className="inline-flex items-center gap-2 text-muted-foreground hover:text-foreground transition-colors"
                  >
                    <Icon className="h-4 w-4" aria-hidden />
                    <span className="font-mono text-[11px] uppercase tracking-[0.12em]">
                      {s.displayLabel}
                    </span>
                  </a>
                </li>
              );
            })}
          </ul>
        </div>
      </div>
    </section>
  );
}
