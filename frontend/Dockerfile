# Local-dev image for the Next.js frontend used by docker-compose.
# Runs `next dev` with hot reload (the compose file bind-mounts the source).
# Production deploys go via Vercel directly from the repo — Vercel ignores this file.

FROM node:20-alpine
WORKDIR /app

RUN apk add --no-cache libc6-compat

COPY package.json package-lock.json ./
RUN npm ci --no-audit --no-fund

COPY . .

ENV NEXT_TELEMETRY_DISABLED=1
EXPOSE 3000

CMD ["npx", "next", "dev", "-H", "0.0.0.0", "-p", "3000"]
