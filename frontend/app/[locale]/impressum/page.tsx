import { setRequestLocale, getTranslations } from "next-intl/server";

export default async function ImpressumPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  setRequestLocale(locale);
  const t = await getTranslations({ locale, namespace: "legal.impressum" });
  return (
    <article className="prose max-w-2xl mx-auto py-24 space-y-6">
      <h1 className="font-display text-3xl md:text-4xl tracking-tighter">
        {t("title")}
      </h1>
      <p className="label-mono">{t("intro")}</p>
      <div className="text-base text-muted-foreground space-y-2">
        <p>Nabeel Thurakkal</p>
        <p>Ulm, Germany</p>
        <p>nabeel.thurakkal92@gmail.com</p>
      </div>
      <p className="text-sm text-muted-foreground italic">{t("comingSoon")}</p>
    </article>
  );
}

export const metadata = { robots: { index: false, follow: true } };
