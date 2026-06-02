# portfolio-claude

Personal portfolio for **Nabeel Thurakkal** — Senior Frontend Developer, Ulm, Germany.

Multilingual (EN / DE), statically generated, content-driven from a Go + Postgres backend.

## Architecture

```
┌───────────────┐  build-time fetch   ┌───────────────┐    pgx     ┌────────────┐
│   Next.js 15  │ ─────────────────▶  │   Go backend  │ ─────────▶ │ Postgres   │
│   App Router  │                     │  chi + pgx    │            │ 16 / 15    │
│   /frontend   │ ◀──── /api/contact ─┤  Resend       │            │            │
└───────────────┘   runtime proxy     └───────────────┘            └────────────┘
```

- **/frontend** — Next.js 15 (App Router), TypeScript, Tailwind v3.4, next-intl, next-themes, react-hook-form + zod. SSG per locale via `generateStaticParams`. Contact form posts to `/api/contact` which proxies to the Go service.
- **/backend** — Go 1.25, chi router, pgx pool, golang-migrate (embedded), Resend HTTP client. All page content is served as one JSON payload per locale from `GET /api/content?locale=…`.
- **Postgres** — typed tables per entity (`hero`, `about`, `projects`, …) with `*_translations` tables for localized fields. Migrations are embedded in the Go binary and run on startup.

> Note: the original brief specified `sqlc + golang-migrate CLI`. We use pgx directly with hand-written SQL and import golang-migrate as a Go library so no extra CLIs are required for setup. Same architecture, less ceremony.

## Run locally — Docker only

**The only prerequisite is Docker Desktop.** No local Go, Node, or Postgres needed.

```sh
docker compose up
```

That's it. Three services come up:

| Service | Image source | Host port | What it does |
|---|---|---|---|
| `postgres` | `postgres:16-alpine` | `5433` | Database (mapped to host `:5433` to avoid colliding with a local Postgres on `:5432`) |
| `backend`  | `./backend/Dockerfile` | `8080` | Go API. Runs embedded migrations + seeds EN+DE content on first boot |
| `frontend` | `./frontend/Dockerfile` | `3000` | Next.js `next dev` with hot reload — `frontend/` is bind-mounted into the container |

Open <http://localhost:3000> — the middleware redirects to `/en`. The German UI is at `/de` (or via the nav language switcher).

Endpoints exposed:

- Frontend: <http://localhost:3000>
- Backend `/healthz`, `/api/content?locale={en|de}`, `/api/contact`: <http://localhost:8080>
- Postgres on the host: `psql -h localhost -p 5433 -U portfolio -d portfolio` (password `portfolio`)

### Common commands

| Goal | Command |
|---|---|
| First start (build + run) | `docker compose up` |
| Run in background | `docker compose up -d` |
| Stop | `docker compose down` |
| Wipe Postgres data too | `docker compose down -v` |
| Rebuild after backend code change | `docker compose up --build backend` |
| Rebuild after backend deps change | `docker compose build --no-cache backend` |
| Stream logs for one service | `docker compose logs -f frontend` |
| Re-seed content without restart | `docker compose exec backend /app/seed` *(after `docker compose build backend` if you edited `seed.go`)* |
| Shell into a container | `docker compose exec backend sh` |

### How hot reload works

- **Frontend**: source is bind-mounted, `next dev` watches files with polling (`WATCHPACK_POLLING=true` for reliable detection inside Docker on macOS). Edit `frontend/components/...` → browser reloads in ~1 s. `node_modules` and `.next` live inside the container, so npm changes don't fight Docker.
- **Backend**: the existing Dockerfile compiles a Go binary, so backend code changes need `docker compose up --build backend`. Acceptable for a portfolio; for tighter iteration use `cd backend && go run ./cmd/server` natively (you'd need Go installed in that case).

### Passing the Resend API key

`backend/.env` is **not** read inside the container. Set it via the host environment when starting compose:

```sh
RESEND_API_KEY="re_xxx" docker compose up
```

Or, for persistent local config, create a `.env` next to `docker-compose.yml`:

```sh
echo 'RESEND_API_KEY=re_xxx' >> .env       # picked up by `docker compose` automatically
docker compose up
```

Without a key, the contact form still validates and stores submissions to `contact_submissions`, but no email is sent (you'll see a warning in the backend logs).

### Backend endpoints

- `GET  /api/content?locale={en|de}` — full content payload (called by Next when pages render)
- `POST /api/contact` — contact form (Resend, rate-limited, honeypot)
- `GET  /healthz` — liveness

## Content editing

All page content lives in Postgres and is shipped by `backend/internal/seed/seed.go`. Edit the values there (EN + DE side by side), restart the backend, and content reloads (the seed wipes + reinserts content tables in one transaction; submissions and rate-limit rows are not touched).

Adding a new locale:
1. Add the locale code to `routing.locales` in `frontend/i18n/routing.ts`.
2. Add a `frontend/messages/<code>.json` UI strings file.
3. Add seed entries for the new locale in `backend/internal/seed/seed.go`.
4. Restart the backend; rebuild the frontend.

## Contact form (Resend)

- Set `RESEND_API_KEY` in `backend/.env`. Without it, submissions are still stored in `contact_submissions`, but no email is sent (the server logs a warning).
- Set `CONTACT_FROM` to a sender on a domain verified in Resend. Default: `no-reply@thurakkal.com`.
- Set `CONTACT_TO` to your inbox. Default: `nabeel.thurakkal92@gmail.com`. The Resend `reply_to` header is set to the submitter's email so replying in your inbox goes back to them.

### DNS records to add at Namecheap for `thurakkal.com`

Resend will display the exact values after you add the domain in their dashboard. Typical records:

| Type | Host | Value |
|------|------|-------|
| TXT  | `@` or `send` | `v=spf1 include:_spf.resend.com ~all` |
| TXT  | `resend._domainkey` | (DKIM value supplied by Resend) |
| TXT  | `_dmarc` | `v=DMARC1; p=none; rua=mailto:dmarc@thurakkal.com` |

You do **not** need MX records unless you also want to receive mail at that domain.

## SEO & a11y notes

- Per-locale title, description, canonical, OG/Twitter card, hreflang alternates, `Person` JSON-LD with GitHub/LinkedIn `sameAs`.
- `sitemap.xml` and `robots.txt` are generated by Next.
- Skip-to-content link, focus rings, semantic landmarks, `aria-invalid` / `aria-describedby` on form fields, focus moves to the first error.
- `prefers-reduced-motion` disables the grayscale → color hover transition on project cards.

## Deployment — Vercel + Fly.io + Neon

Free tier across the stack, ~10 minutes end-to-end. **Run in this order** — the Vercel build fetches content from the live backend, so the backend has to be up first.

### 1. Database — Neon

1. Create a project at <https://neon.tech> (free tier, EU region for low latency from Fly Frankfurt).
2. Copy the **pooled** connection string (`...pooler.neon.tech...`) — looks like `postgres://user:pass@ep-xxx-pooler.eu-central-1.aws.neon.tech/neondb?sslmode=require`.

### 2. Backend — Fly.io

Install once: `brew install flyctl && fly auth signup`.

```sh
cd backend
fly launch --copy-config --no-deploy        # registers the app using ./fly.toml
fly secrets set \
  DATABASE_URL="postgres://...pooler.neon.tech/...?sslmode=require" \
  RESEND_API_KEY="re_xxx"                   # leave blank for now if you haven't verified yet
fly deploy
```

The first deploy runs embedded migrations and seeds EN+DE content. Verify:
```sh
curl https://thurakkal-portfolio-backend.fly.dev/healthz   # → 204
curl https://thurakkal-portfolio-backend.fly.dev/api/content?locale=en | head -c 200
```

Subsequent backend changes deploy automatically via `.github/workflows/deploy-backend.yml` — just push to `main`. Set `FLY_API_TOKEN` once in **GitHub repo → Settings → Secrets and variables → Actions**: get the token with `fly tokens create deploy -x 999999h`.

### 3. Frontend — Vercel

In the Vercel dashboard:
1. **New Project** → import `thurakkal92/portfolio-claude`.
2. **Root Directory** = `frontend`. Vercel auto-detects Next.js and runs `npm install` + `next build` from there.
3. **Environment Variables**:
   - `API_BASE_URL` = `https://thurakkal-portfolio-backend.fly.dev`
   - `NEXT_PUBLIC_SITE_URL` = `https://thurakkal.com`
4. Deploy.

After it succeeds: **Project → Settings → Domains → Add** `thurakkal.com` and `www.thurakkal.com`. Vercel shows the exact DNS records to set.

### 4. DNS — Namecheap (`thurakkal.com`)

In **Namecheap → Advanced DNS**, replace the default records with:

| Type | Host | Value | TTL |
|---|---|---|---|
| A | `@` | `76.76.21.21` | Automatic |
| CNAME | `www` | `cname.vercel-dns.com.` | Automatic |
| CNAME | `api` | `thurakkal-portfolio-backend.fly.dev.` | Automatic |

Plus the Resend records once you add `thurakkal.com` in the Resend dashboard:

| Type | Host | Value |
|---|---|---|
| TXT  | `send` | `v=spf1 include:_spf.resend.com ~all` |
| TXT  | `resend._domainkey` | *(DKIM value Resend gives you)* |
| TXT  | `_dmarc` | `v=DMARC1; p=none; rua=mailto:dmarc@thurakkal.com` |

Propagation: usually 5–15 minutes on Namecheap.

### 5. Optional — custom backend hostname

If you want the backend at `https://api.thurakkal.com` instead of the `.fly.dev` URL:
```sh
cd backend
fly certs add api.thurakkal.com
```
Then update `API_BASE_URL` in Vercel and `ALLOWED_ORIGINS` in `fly.toml` accordingly, and redeploy.

### Cost summary

| Service | Plan | What's included | Cost |
|---|---|---|---|
| Vercel | Hobby | 100 GB-hours, automatic deploys, custom domain | $0 |
| Fly.io | Free | 3× shared-cpu-1x 256 MB VMs, 160 GB egress | $0 |
| Neon | Free | 0.5 GB storage, autoscale to zero, branching | $0 |
| Resend | Free | 3,000 emails/mo, 100/day | $0 |
| **Total** | | | **$0/mo** |

## Layout

```
backend/
  cmd/server, cmd/seed
  internal/{config,db,content,contact,http,i18n,seed}
  migrations/*.sql            # embedded in the binary
  embed.go                    # exposes migrations FS
frontend/
  app/[locale]/{layout,page,impressum,datenschutz}
  app/api/contact/route.ts    # proxy to backend
  app/{sitemap,robots,not-found}.ts(x)
  components/{ui,sections,nav,footer,contact-form,theme-*,language-switcher}
  i18n/{routing,request}.ts
  lib/{api,types,seo,utils}.ts
  messages/{en,de}.json
  public/{cv,images/projects,og}
docker-compose.yml            # `docker compose up` runs all three services
```

## Replace before launch

- Project screenshots in `frontend/public/images/projects/` (currently SVG placeholders).
- OG images in `frontend/public/og/` (currently unset; metadata references `og-image-{en,de}.png`).
- Impressum + Datenschutz copy (currently placeholders; legally required for a DE-resident operator).
- Resend domain verification + DKIM records at Namecheap.
