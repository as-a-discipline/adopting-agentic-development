# Multi-stage build for the Pulse web app. Used by the Compose integration
# environment (deploy/compose/docker-compose.yml); day-to-day development uses
# the containerized `task web:*` targets instead (see web/Taskfile.yml and
# web/AGENTS.md) — this Dockerfile packages the same build.mjs output behind a
# static file server for Compose/Helm.
FROM node:22-alpine AS build
WORKDIR /app
COPY . .
RUN node tools/build.mjs

# Unprivileged variant: runs as a non-root user (see policies/README.md rule
# 5, container/deployment quality) rather than the stock nginx:*-alpine image,
# which runs its master process as root.
FROM nginxinc/nginx-unprivileged:1.27.5-alpine
COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
# Explicit even though the base image already defaults to this user — makes
# the non-root runtime user verifiable directly from this Dockerfile (see
# scripts/check-container-quality.sh / policies/README.md rule 5) without
# needing to inspect the upstream base image's own Dockerfile.
USER nginx
EXPOSE 8080
