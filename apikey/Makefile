# Makefile
SHELL := /bin/bash
.PHONY: help test clean deps install-tools check-go-version fmt vet lint gosec

REPORTS_DIR := reports
GO_VERSION ?= 1.26
GOLANGCI_LINT_VERSION ?= 2.10
GOSEC_VERSION ?= 2.23

help: ## Mostrar comandos disponibles
	@echo "Connect-Libraries - Comandos de desarrollo"
	@echo "=========================================="
	@echo "  install-tools instala golangci-lint y gosec en versiones definidas"
	@echo "  check-go-version valida version minima de Go requerida"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ { printf "  %-12s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

test: ## Ejecutar tests
	@echo "Ejecutando tests..."
	@go test -v ./...

clean: ## Limpiar cache y reportes
	@echo "Limpiando..."
	@go clean -cache -testcache -modcache
	@rm -rf $(REPORTS_DIR)

deps: ## Descargar y ordenar dependencias
	@echo "Actualizando dependencias..."
	@go mod download
	@go mod tidy

check-go-version: ## Validar version minima de Go
	@echo "Validando version de Go..."
	@installed=$$(go version 2>/dev/null | sed -nE 's/^go version go([0-9]+(\.[0-9]+){1,2}).*/\1/p'); \
	[ -n "$$installed" ] || { echo "Error: no se pudo detectar la version de Go instalada."; exit 1; }; \
	required="$(GO_VERSION)"; required=$${required#v}; \
	case "$$installed" in \
		$$required|$$required.*) ;; \
		*) echo "Error: se requiere Go $$required.x (instalada $$installed)."; exit 1 ;; \
	esac

install-tools: check-go-version ## Instalar herramientas de CI/desarrollo
	@echo "Instalando herramientas de desarrollo..."
	@set -euo pipefail; \
	lint_version="$(GOLANGCI_LINT_VERSION)"; lint_version=$${lint_version#v}; \
	gosec_version="$(GOSEC_VERSION)"; gosec_version=$${gosec_version#v}; \
	go install "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v$${lint_version}"; \
	go install "github.com/securego/gosec/v2/cmd/gosec@v$${gosec_version}"; \
	echo "Herramientas instaladas: golangci-lint v$${lint_version}, gosec v$${gosec_version}"

fmt: ## Formatear codigo
	@echo "Formateando codigo..."
	@go fmt ./...

vet: ## Ejecutar go vet
	@echo "Ejecutando go vet..."
	@go vet ./...

lint: ## Ejecutar golangci-lint
	@echo "Ejecutando golangci-lint..."
	@command -v golangci-lint >/dev/null 2>&1 || { echo "Error: golangci-lint no esta instalado. Instalalo con la version $(GOLANGCI_LINT_VERSION)."; exit 1; }
	@installed=$$(golangci-lint version 2>/dev/null | head -n1); \
	required="$(GOLANGCI_LINT_VERSION)"; required=$${required#v}; \
	installed_version=$$(echo "$$installed" | sed -nE 's/.*version[[:space:]]+v?([0-9]+(\.[0-9]+){1,2}).*/\1/p' | head -n1); \
	[ -n "$$installed_version" ] || { echo "Error: no se pudo detectar la version instalada de golangci-lint."; exit 1; }; \
	case "$$installed_version" in \
		$$required|$$required.*) ;; \
		*) echo "Error: se requiere golangci-lint $$required.x (instalada $$installed_version)."; exit 1 ;; \
	esac; \
	echo "$$installed" | grep -q "built with go$(GO_VERSION)" || { echo "Error: se requiere golangci-lint compilado con go$(GO_VERSION)."; exit 1; }
	@mkdir -p $(REPORTS_DIR)
	@if golangci-lint run --help 2>&1 | grep -q -- "--output.json.path"; then \
		golangci-lint run --timeout=5m --output.json.path $(REPORTS_DIR)/lint.json; \
	else \
		golangci-lint run --timeout=5m --out-format json > $(REPORTS_DIR)/lint.json; \
	fi
	@echo "Reporte JSON generado en $(REPORTS_DIR)/lint.json"

gosec: ## Ejecutar escaneo de seguridad con gosec (mismo alcance que CI)
	@echo "Ejecutando escaneo de seguridad con gosec..."
	@command -v gosec >/dev/null 2>&1 || { echo "Error: gosec no esta instalado. Instalalo con la version $(GOSEC_VERSION)."; exit 1; }
	@installed=$$(gosec -version 2>&1); \
	required="$(GOSEC_VERSION)"; required=$${required#v}; \
	installed_version=$$(echo "$$installed" | grep -Eo 'v?[0-9]+\.[0-9]+(\.[0-9]+)?' | head -n1 | sed 's/^v//'); \
	if [ -z "$$installed_version" ]; then \
		gosec_bin=$$(command -v gosec || true); \
		if [ -n "$$gosec_bin" ]; then \
			module_line=$$(go version -m "$$gosec_bin" 2>/dev/null | grep -E '^mod[[:space:]]+github.com/securego/gosec/v2[[:space:]]+v' | head -n1 || true); \
			installed_version=$$(echo "$$module_line" | grep -Eo 'v[0-9]+\.[0-9]+(\.[0-9]+)?' | head -n1 | sed 's/^v//'); \
		fi; \
	fi; \
	if [ -z "$$installed_version" ]; then \
		gosec_bin_gopath="$$(go env GOPATH)/bin/gosec"; \
		if [ -x "$$gosec_bin_gopath" ]; then \
			module_line=$$(go version -m "$$gosec_bin_gopath" 2>/dev/null | grep -E '^mod[[:space:]]+github.com/securego/gosec/v2[[:space:]]+v' | head -n1 || true); \
			installed_version=$$(echo "$$module_line" | grep -Eo 'v[0-9]+\.[0-9]+(\.[0-9]+)?' | head -n1 | sed 's/^v//'); \
		fi; \
	fi; \
	if [ -n "$$installed_version" ]; then \
		case "$$installed_version" in \
			$$required|$$required.*) ;; \
			*) echo "Error: se requiere gosec $$required.x (instalada $$installed_version)."; exit 1 ;; \
		esac; \
	else \
		echo "Aviso: no se pudo detectar version numerica de gosec; se omite validacion de version."; \
	fi

	@mkdir -p $(REPORTS_DIR)
	@rc=0; \
	gosec -fmt=json -out=$(REPORTS_DIR)/gosec.json ./... || rc=$$?; \
	echo "Reporte JSON generado en $(REPORTS_DIR)/gosec.json"; \
	if command -v jq >/dev/null 2>&1; then \
		found=$$(jq -r '.Stats.found // 0' $(REPORTS_DIR)/gosec.json); \
		files=$$(jq -r '.Stats.files // 0' $(REPORTS_DIR)/gosec.json); \
		nosec=$$(jq -r '.Stats.nosec // 0' $(REPORTS_DIR)/gosec.json); \
		printf "Resumen gosec: found=%s files=%s nosec=%s\n" "$$found" "$$files" "$$nosec"; \
	else \
		found=$$(grep -o '"rule_id"' $(REPORTS_DIR)/gosec.json | wc -l | tr -d '[:space:]'); \
		printf "Resumen gosec: found=%s (instala jq para ver mas detalle)\n" "$$found"; \
	fi; \
	if [ $$rc -ne 0 ]; then \
		echo "gosec finalizo con errores (exit $$rc). Revisa $(REPORTS_DIR)/gosec.json"; \
		exit $$rc; \
	fi
