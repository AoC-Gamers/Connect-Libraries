# Makefile
SHELL := /bin/bash
.PHONY: help report test lint gosec deps fmt vet clean clear install-tools check-go-version

LIBRARIES := apikey audit authz errors middleware migrate nats swagger testhelpers
GO_VERSION ?= 1.26
GOLANGCI_LINT_VERSION ?= 2.10
GOSEC_VERSION ?= v2.23.0

help: ## Mostrar comandos disponibles
	@echo "Connect-Libraries - Reportes por subdirectorio"
	@echo "================================================"
	@echo "  install-tools instala golangci-lint y gosec en todas las librerias"
	@echo "  check-go-version valida version minima de Go requerida"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ { printf "  %-16s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

report: test lint gosec ## Ejecutar test, lint y gosec para todas las librerias

check-go-version: ## Validar version minima de Go
	@echo "Validando version de Go..."
	@installed=$$(go version 2>/dev/null | sed -nE 's/^go version go([0-9]+(\.[0-9]+){1,2}).*/\1/p'); \
	[ -n "$$installed" ] || { echo "Error: no se pudo detectar la version de Go instalada."; exit 1; }; \
	required="$(GO_VERSION)"; required=$${required#v}; \
	case "$$installed" in \
		$$required|$$required.*) ;; \
		*) echo "Error: se requiere Go $$required.x (instalada $$installed)."; exit 1 ;; \
	esac

install-tools: check-go-version ## Instalar herramientas de CI/desarrollo en todas las librerias
	@status=0; \
	for lib in $(LIBRARIES); do \
		echo "[TOOLS] $$lib"; \
		if $(MAKE) -C $$lib install-tools GO_VERSION=$(GO_VERSION) GOLANGCI_LINT_VERSION=$(GOLANGCI_LINT_VERSION) GOSEC_VERSION=$(GOSEC_VERSION); then \
			true; \
		else \
			status=1; \
		fi; \
	done; \
	exit $$status

deps: ## Descargar y ordenar dependencias en todas las librerias
	@status=0; \
	for lib in $(LIBRARIES); do \
		echo "[DEPS] $$lib"; \
		if $(MAKE) -C $$lib deps; then \
			true; \
		else \
			status=1; \
		fi; \
	done; \
	exit $$status

fmt: ## Formatear codigo en todas las librerias
	@status=0; \
	for lib in $(LIBRARIES); do \
		echo "[FMT] $$lib"; \
		if $(MAKE) -C $$lib fmt; then \
			true; \
		else \
			status=1; \
		fi; \
	done; \
	exit $$status

vet: ## Ejecutar go vet en todas las librerias
	@status=0; \
	for lib in $(LIBRARIES); do \
		echo "[VET] $$lib"; \
		if $(MAKE) -C $$lib vet; then \
			true; \
		else \
			status=1; \
		fi; \
	done; \
	exit $$status

clean: ## Limpiar cache y reportes en todas las librerias
	@status=0; \
	for lib in $(LIBRARIES); do \
		echo "[CLEAN] $$lib"; \
		if $(MAKE) -C $$lib clean; then \
			true; \
		else \
			status=1; \
		fi; \
	done; \
	exit $$status

test: ## Generar reportes de test en <libreria>/reports/test.log
	@status=0; \
	for lib in $(LIBRARIES); do \
		lib_reports="$$lib/reports"; \
		echo "[TEST] $$lib"; \
		mkdir -p "$$lib_reports"; \
		if $(MAKE) -C $$lib test > "$$lib_reports/test.log" 2>&1; then \
			echo "  OK  -> $$lib_reports/test.log"; \
		else \
			echo "  FAIL -> $$lib_reports/test.log"; \
			status=1; \
		fi; \
	done; \
	exit $$status

lint: ## Generar reportes de lint en <libreria>/reports/lint.json y lint.log
	@status=0; \
	for lib in $(LIBRARIES); do \
		lib_reports="$$lib/reports"; \
		echo "[LINT] $$lib"; \
		mkdir -p "$$lib_reports"; \
		if $(MAKE) -C $$lib lint GO_VERSION=$(GO_VERSION) GOLANGCI_LINT_VERSION=$(GOLANGCI_LINT_VERSION); then \
			echo "  OK  -> $$lib_reports/lint.json"; \
		else \
			echo "  FAIL -> $$lib_reports/lint.json"; \
			status=1; \
		fi; \
	done; \
	exit $$status

gosec: ## Reproducir gosec de CI para todas las librerias (logs en <libreria>/reports/gosec.log)
	@status=0; \
	for lib in $(LIBRARIES); do \
		lib_reports="$$lib/reports"; \
		echo "[GOSEC] $$lib"; \
		mkdir -p "$$lib_reports"; \
		if $(MAKE) -C $$lib gosec GOSEC_VERSION=$(GOSEC_VERSION); then \
			echo "  OK  -> $$lib_reports/gosec.log"; \
		else \
			echo "  FAIL -> $$lib_reports/gosec.log"; \
			status=1; \
		fi; \
	done; \
	exit $$status

clear: ## Eliminar reportes en cada libreria
	@for lib in $(LIBRARIES); do \
		rm -rf "$$lib/reports"; \
		echo "Reportes eliminados: $$lib/reports"; \
	done
