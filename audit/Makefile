# Makefile
.PHONY: help test clean deps fmt vet lint gosec

REPORTS_DIR := reports
GO_VERSION ?= 1.26
GOLANGCI_LINT_VERSION ?= 2.10
GOSEC_VERSION ?= v2.23.0

help: ## Mostrar comandos disponibles
	@echo "Connect-Libraries - Comandos de desarrollo"
	@echo "=========================================="
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

gosec: ## Ejecutar escaneo de seguridad con gosec
	@echo "Ejecutando escaneo de seguridad con gosec (nativo)..."
	@mkdir -p $(REPORTS_DIR)
	@command -v gosec >/dev/null 2>&1 || { echo "Error: gosec no esta instalado. Instalalo con la version $(GOSEC_VERSION)."; exit 1; }
	@installed=$$(gosec -version 2>&1); \
	required="$(GOSEC_VERSION)"; required=$${required#v}; \
	installed_version=$$(echo "$$installed" | grep -Eo 'v?[0-9]+\.[0-9]+(\.[0-9]+)?' | head -n1 | sed 's/^v//'); \
	if [ -n "$$installed_version" ]; then \
		case "$$installed_version" in \
			$$required|$$required.*) ;; \
			*) echo "Error: se requiere gosec $$required.x (instalada $$installed_version)."; exit 1 ;; \
		esac; \
	else \
		echo "Aviso: no se pudo detectar version numerica de gosec; se omite validacion de version."; \
	fi
	@rc=0; \
	NO_COLOR=1 gosec ./... 2>&1 | sed -E 's/\x1b\[[0-9;]*m//g' | tee $(REPORTS_DIR)/gosec.log || rc=$$?; \
	if [ $$rc -ne 0 ]; then \
		echo "gosec finalizo con errores (exit $$rc). Revisa $(REPORTS_DIR)/gosec.log"; \
		exit $$rc; \
	fi; \
	echo "Reporte generado en $(REPORTS_DIR)/gosec.log"
