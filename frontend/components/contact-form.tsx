"use client";

import * as React from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useTranslations } from "next-intl";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import type { ContactCopy } from "@/lib/types";

type Props = {
  copy: ContactCopy;
  locale: string;
};

type Values = {
  name: string;
  email: string;
  message: string;
  website: string; // honeypot
};

type Status = { kind: "idle" } | { kind: "success" } | { kind: "error"; message: string };

export function ContactForm({ copy, locale }: Props) {
  const t = useTranslations("form");

  const schema = React.useMemo(
    () =>
      z.object({
        name: z.string().min(1, t("errors.required")),
        email: z
          .string()
          .min(1, t("errors.required"))
          .email(t("errors.email")),
        message: z
          .string()
          .min(1, t("errors.required"))
          .min(10, t("errors.minMessage")),
        website: z.string().optional().default(""),
      }),
    [t],
  );

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    setError,
    setFocus,
    reset,
  } = useForm<Values>({
    resolver: zodResolver(schema),
    defaultValues: { name: "", email: "", message: "", website: "" },
  });

  const [status, setStatus] = React.useState<Status>({ kind: "idle" });

  React.useEffect(() => {
    const first =
      (errors.name && "name") ||
      (errors.email && "email") ||
      (errors.message && "message") ||
      null;
    if (first) setFocus(first as keyof Values);
  }, [errors, setFocus]);

  async function onSubmit(values: Values) {
    setStatus({ kind: "idle" });
    try {
      const res = await fetch("/api/contact", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ...values, locale }),
      });
      if (res.ok) {
        reset();
        setStatus({ kind: "success" });
        return;
      }
      if (res.status === 429) {
        setStatus({ kind: "error", message: t("errors.rateLimited") });
        return;
      }
      if (res.status === 422) {
        const body = (await res.json()) as { fields?: Record<string, string> };
        if (body.fields) {
          for (const [k, v] of Object.entries(body.fields)) {
            setError(k as keyof Values, { type: "server", message: v });
          }
          return;
        }
      }
      setStatus({ kind: "error", message: copy.errorMessage });
    } catch {
      setStatus({ kind: "error", message: t("errors.network") });
    }
  }

  return (
    <form
      onSubmit={handleSubmit(onSubmit)}
      noValidate
      className="space-y-4 text-left"
      aria-describedby={status.kind === "error" ? "form-status" : undefined}
    >
      {/* Honeypot */}
      <div className="absolute -left-[9999px]" aria-hidden>
        <label htmlFor="website">Website</label>
        <input
          id="website"
          type="text"
          tabIndex={-1}
          autoComplete="off"
          {...register("website")}
        />
      </div>

      <Field
        id="name"
        label={copy.formNameLabel}
        error={errors.name?.message}
      >
        <Input
          id="name"
          type="text"
          autoComplete="name"
          placeholder={t("namePlaceholder")}
          invalid={!!errors.name}
          aria-describedby={errors.name ? "name-error" : undefined}
          {...register("name")}
        />
      </Field>

      <Field
        id="email"
        label={copy.formEmailLabel}
        error={errors.email?.message}
      >
        <Input
          id="email"
          type="email"
          autoComplete="email"
          placeholder={t("emailPlaceholder")}
          invalid={!!errors.email}
          aria-describedby={errors.email ? "email-error" : undefined}
          {...register("email")}
        />
      </Field>

      <Field
        id="message"
        label={copy.formMessageLabel}
        error={errors.message?.message}
      >
        <Textarea
          id="message"
          rows={4}
          placeholder={t("messagePlaceholder")}
          invalid={!!errors.message}
          aria-describedby={errors.message ? "message-error" : undefined}
          {...register("message")}
        />
      </Field>

      <Button
        type="submit"
        variant="primary"
        size="lg"
        className="w-full"
        disabled={isSubmitting}
      >
        {isSubmitting ? t("sending") : copy.formSubmitLabel}
      </Button>

      <div
        id="form-status"
        role="status"
        aria-live="polite"
        className="text-sm min-h-[1.25rem]"
      >
        {status.kind === "success" ? (
          <span className="text-emerald-600 dark:text-emerald-400">
            {copy.successMessage}
          </span>
        ) : status.kind === "error" ? (
          <span className="text-red-600 dark:text-red-400">{status.message}</span>
        ) : null}
      </div>
    </form>
  );
}

function Field({
  id,
  label,
  error,
  children,
}: {
  id: string;
  label: string;
  error?: string;
  children: React.ReactNode;
}) {
  return (
    <div className="space-y-1.5">
      <label htmlFor={id} className="label-mono">
        {label}
      </label>
      {children}
      {error ? (
        <p id={`${id}-error`} className="text-xs text-red-600 dark:text-red-400">
          {error}
        </p>
      ) : null}
    </div>
  );
}
