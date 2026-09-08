# Multi-stage build for the Pulse Go service. Used by the Compose integration
# environment (deploy/compose/docker-compose.yml); day-to-day development uses
# the containerized `task service:*` targets instead (see service/Taskfile.yml
# and /adr/0003-containerized-tooling.md) — this Dockerfile is the packaging
# artifact for running the built binary as its own container image.
FROM golang:1.23-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/pulse ./cmd/pulse

FROM alpine:3.20
RUN adduser -D -u 10001 pulse
COPY --from=build /out/pulse /usr/local/bin/pulse
USER pulse
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/pulse"]
