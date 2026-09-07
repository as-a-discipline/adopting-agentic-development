# Makefile
#
# This is NOT the project build system. It only handles host prerequisite
# checks and bootstrap. Actual build/test/lint/generate work happens through
# `task` targets (see Taskfile.yml) once Task is installed.

SHELL := /bin/bash

.PHONY: help check-prereqs bootstrap

help: ## Show this help
	@echo "Pulse — host bootstrap"
	@echo ""
	@echo "  make check-prereqs   Verify required host tools are installed"
	@echo "  make bootstrap       Install Task if missing, then check prerequisites"
	@echo ""
	@echo "All project work (build/test/lint/generate/validate) is done via 'task'."
	@echo "Run 'task --list' after bootstrap to see available targets."

check-prereqs: ## Verify required host tools (git, docker, task) are available
	@status=0; \
	command -v git >/dev/null 2>&1 || { echo "MISSING: git — install from https://git-scm.com/downloads"; status=1; }; \
	command -v docker >/dev/null 2>&1 || { echo "MISSING: docker — install Docker Desktop or Docker Engine: https://docs.docker.com/get-docker/"; status=1; }; \
	if command -v docker >/dev/null 2>&1; then \
		docker info >/dev/null 2>&1 || { echo "DOCKER NOT USABLE: 'docker info' failed — is the Docker daemon running?"; status=1; }; \
	fi; \
	command -v task >/dev/null 2>&1 || { echo "MISSING: task — run 'make bootstrap' or see https://taskfile.dev/installation/"; status=1; }; \
	if [ $$status -eq 0 ]; then echo "All required host prerequisites are present."; fi; \
	exit $$status

bootstrap: ## Install Task if missing (best-effort, portable), then verify prerequisites
	@if command -v task >/dev/null 2>&1; then \
		echo "Task already installed: $$(task --version)"; \
	elif command -v brew >/dev/null 2>&1; then \
		echo "Installing Task via Homebrew..."; \
		brew install go-task/tap/go-task; \
	elif command -v go >/dev/null 2>&1; then \
		echo "Installing Task via 'go install'..."; \
		go install github.com/go-task/task/v3/cmd/task@latest; \
	else \
		echo "Could not auto-install Task safely on this host."; \
		echo "Install it manually: https://taskfile.dev/installation/"; \
	fi
	@$(MAKE) --no-print-directory check-prereqs
