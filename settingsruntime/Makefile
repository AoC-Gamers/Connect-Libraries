# Makefile
SHELL := /bin/bash
.PHONY: help test clean deps go-tools fmt vet lint gosec

MODULE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
REPORTS_DIR := $(MODULE_DIR)reports
GOLANGCI_LINT_VERSION ?= 2.10
GOSEC_VERSION ?= 2.23
SCRIPTS_DIR := $(abspath $(MODULE_DIR)/../scripts)

help: ## Mostrar comandos disponibles
	@echo "Connect-Libraries/settingsruntime - Comandos de desarrollo"
	@echo "==========================================================="
	@echo "  go-tools instala golangci-lint y gosec en versiones definidas"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ { printf "  %-14s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

test: ## Ejecutar tests
	@echo "Ejecutando tests..."
	@cd $(MODULE_DIR) && go test -v ./...

clean: ## Limpiar cache y reportes
	@echo "Limpiando..."
	@cd $(MODULE_DIR) && go clean -cache -testcache -modcache
	@rm -rf $(REPORTS_DIR)

deps: ## Descargar y ordenar dependencias
	@echo "Actualizando dependencias..."
	@cd $(MODULE_DIR) && go mod download
	@cd $(MODULE_DIR) && go mod tidy
	@cd $(MODULE_DIR) && go mod verify

go-tools: ## Instalar herramientas de CI/desarrollo
	@bash $(SCRIPTS_DIR)/make-go-tools.sh "$(GOLANGCI_LINT_VERSION)" "$(GOSEC_VERSION)"

fmt: ## Formatear codigo
	@echo "Formateando codigo..."
	@cd $(MODULE_DIR) && go fmt ./...

vet: ## Ejecutar go vet
	@echo "Ejecutando go vet..."
	@cd $(MODULE_DIR) && go vet ./...

lint: ## Ejecutar golangci-lint
	@cd $(MODULE_DIR) && bash $(SCRIPTS_DIR)/make-lint.sh "$(REPORTS_DIR)" "$(GOLANGCI_LINT_VERSION)"

gosec: ## Ejecutar escaneo de seguridad con gosec (mismo alcance que CI)
	@cd $(MODULE_DIR) && bash $(SCRIPTS_DIR)/make-gosec.sh "$(REPORTS_DIR)" "$(GOSEC_VERSION)"
