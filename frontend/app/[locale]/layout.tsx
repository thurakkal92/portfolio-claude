import { notFound } from "next/navigation";
import { NextIntlClientProvider } from "next-intl";
import { getMessages, getTranslations, setRequestLocale } from "next-intl/server";
import { Geist, Inter, JetBrains_Mono } from "next/font/google";
import { ThemeProvider } from "@/components/theme-provider";
import { Nav } from "@/components/nav";
import { Footer } from "@/components/footer";
import { routing } from "@/i18n/routing";

const geist = Geist({
  subsets: ["latin"],
  display: "swap",
  variable: "--font-display",
});
const inter = Inter({
  subsets: ["latin"],
  display: "swap",
  variable: "--font-body",
});
const jetbrains = JetBrains_Mono({
  subsets: ["latin"],
  display: "swap",
  variable: "--font-mono",
});

export function generateStaticParams() {
  return routing.locales.map((locale) => ({ locale }));
}

export default async function LocaleLayout({
  children,
  params,
}: {
  children: React.ReactNode;
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  if (!routing.locales.includes(locale as (typeof routing.locales)[number])) {
    notFound();
  }
  setRequestLocale(locale);
  const messages = await getMessages();
  const nav = await getTranslations({ locale, namespace: "nav" });

  return (
    <html
      lang={locale}
      suppressHydrationWarning
      className={`${geist.variable} ${inter.variable} ${jetbrains.variable}`}
    >
      <body>
        <NextIntlClientProvider messages={messages} locale={locale}>
          <ThemeProvider>
            <a
              href="#main"
              className="sr-only focus:not-sr-only focus:absolute focus:top-4 focus:left-4 focus:z-50 focus:bg-foreground focus:text-background focus:px-3 focus:py-2 focus:font-mono focus:text-xs focus:uppercase"
            >
              {nav("skipToContent")}
            </a>
            <Nav locale={locale} />
            <main id="main" className="container mx-auto">
              {children}
            </main>
            <Footer locale={locale} />
          </ThemeProvider>
        </NextIntlClientProvider>
      </body>
    </html>
  );
}
