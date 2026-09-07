# Multi-stage build for the Pulse web app. Used by the Compose integration
# environment (deploy/compose/docker-compose.yml); day-to-day development uses
# the containerized `task web:*` targets instead (see web/Taskfile.yml and
# web/AGENTS.md) — this Dockerfile packages the same build.mjs output behind a
# static file server for Compose/Helm.
FROM node:22-alpine AS build
WORKDIR /app
COPY . .
RUN node tools/build.mjs

FROM nginx:1.27.5-alpine
COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
