# Makefile
.PHONY: help test clean deps fmt vet lint gosec

REPORTS_DIR := reports
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
	@command -v golangci-lint >/dev/null 2>&1 || { echo "Error: golangci-lint no esta instalado. Instalalo con: scoop install golangci-lint"; exit 1; }
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
	@command -v gosec >/dev/null 2>&1 || { \
		echo "Error: gosec no esta instalado."; \
		echo "Instala gosec version $(GOSEC_VERSION) con tu gestor de paquetes preferido."; \
		exit 1; \
	}
	@rc=0; \
	NO_COLOR=1 gosec ./... 2>&1 | sed -E 's/\x1b\[[0-9;]*m//g' | tee $(REPORTS_DIR)/gosec.log || rc=$$?; \
	if [ $$rc -ne 0 ]; then \
		echo "gosec finalizo con errores (exit $$rc). Revisa $(REPORTS_DIR)/gosec.log"; \
		exit $$rc; \
	fi; \
	echo "Reporte generado en $(REPORTS_DIR)/gosec.log"
