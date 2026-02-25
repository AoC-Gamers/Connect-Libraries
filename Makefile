# Makefile
.PHONY: help report report-test report-lint report-security clean-reports

REPORTS_DIR ?= reports
LIBRARIES := apikey audit authz errors middleware migrate nats swagger testhelpers

help: ## Mostrar comandos disponibles
	@echo "Connect-Libraries - Reportes por subdirectorio"
	@echo "================================================"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ { printf "  %-16s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

report: report-test report-lint report-security ## Ejecutar test, lint y gosec para todas las librerias

report-test: ## Generar reportes de test en reports/<libreria>/test.log
	@mkdir -p $(REPORTS_DIR)
	@status=0; \
	for lib in $(LIBRARIES); do \
		echo "[TEST] $$lib"; \
		mkdir -p $(REPORTS_DIR)/$$lib; \
		if $(MAKE) -C $$lib test > $(REPORTS_DIR)/$$lib/test.log 2>&1; then \
			echo "  OK  -> $(REPORTS_DIR)/$$lib/test.log"; \
		else \
			echo "  FAIL -> $(REPORTS_DIR)/$$lib/test.log"; \
			status=1; \
		fi; \
	done; \
	exit $$status

report-lint: ## Generar reportes de lint en reports/<libreria>/lint.json y lint.log
	@mkdir -p $(REPORTS_DIR)
	@status=0; \
	for lib in $(LIBRARIES); do \
		echo "[LINT] $$lib"; \
		mkdir -p $(REPORTS_DIR)/$$lib; \
		if $(MAKE) -C $$lib lint REPORTS_DIR=../$(REPORTS_DIR)/$$lib > $(REPORTS_DIR)/$$lib/lint.log 2>&1; then \
			echo "  OK  -> $(REPORTS_DIR)/$$lib/lint.json"; \
		else \
			echo "  FAIL -> $(REPORTS_DIR)/$$lib/lint.log"; \
			status=1; \
		fi; \
	done; \
	exit $$status

report-security: ## Generar reportes de seguridad en reports/<libreria>/gosec.json y gosec.log
	@mkdir -p $(REPORTS_DIR)
	@status=0; \
	for lib in $(LIBRARIES); do \
		echo "[GOSEC] $$lib"; \
		mkdir -p $(REPORTS_DIR)/$$lib; \
		if $(MAKE) -C $$lib gosec REPORTS_DIR=../$(REPORTS_DIR)/$$lib > $(REPORTS_DIR)/$$lib/gosec.log 2>&1; then \
			echo "  OK  -> $(REPORTS_DIR)/$$lib/gosec.json"; \
		else \
			echo "  FAIL -> $(REPORTS_DIR)/$$lib/gosec.log"; \
			status=1; \
		fi; \
	done; \
	exit $$status

clean-reports: ## Eliminar reportes agregados de la raiz
	@rm -rf $(REPORTS_DIR)
	@echo "Reportes eliminados: $(REPORTS_DIR)"
