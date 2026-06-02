export default function NotFound() {
  return (
    <html lang="en">
      <body className="min-h-screen flex items-center justify-center">
        <div className="text-center space-y-3">
          <p className="font-mono text-xs uppercase tracking-[0.12em]">404</p>
          <h1 className="font-display text-3xl">Page not found</h1>
          <a
            href="/en"
            className="inline-block font-mono text-xs uppercase tracking-[0.08em] underline underline-offset-4"
          >
            Back home
          </a>
        </div>
      </body>
    </html>
  );
}
