# Makefile
.PHONY: help report report-test report-lint report-security clean-reports

LIBRARIES := apikey audit authz errors middleware migrate nats swagger testhelpers

help: ## Mostrar comandos disponibles
	@echo "Connect-Libraries - Reportes por subdirectorio"
	@echo "================================================"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ { printf "  %-16s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

report: report-test report-lint report-security ## Ejecutar test, lint y gosec para todas las librerias

report-test: ## Generar reportes de test en <libreria>/reports/test.log
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

report-lint: ## Generar reportes de lint en <libreria>/reports/lint.json y lint.log
	@status=0; \
	for lib in $(LIBRARIES); do \
		lib_reports="$$lib/reports"; \
		echo "[LINT] $$lib"; \
		mkdir -p "$$lib_reports"; \
		if $(MAKE) -C $$lib lint > "$$lib_reports/lint.log" 2>&1; then \
			echo "  OK  -> $$lib_reports/lint.json"; \
		else \
			echo "  FAIL -> $$lib_reports/lint.log"; \
			status=1; \
		fi; \
	done; \
	exit $$status

report-security: ## Generar reportes de seguridad en <libreria>/reports/gosec.json y gosec.log
	@status=0; \
	for lib in $(LIBRARIES); do \
		lib_reports="$$lib/reports"; \
		echo "[GOSEC] $$lib"; \
		mkdir -p "$$lib_reports"; \
		if $(MAKE) -C $$lib gosec > "$$lib_reports/gosec.log" 2>&1; then \
			echo "  OK  -> $$lib_reports/gosec.json"; \
		else \
			echo "  FAIL -> $$lib_reports/gosec.log"; \
			status=1; \
		fi; \
	done; \
	exit $$status

clean-reports: ## Eliminar reportes en cada libreria
	@for lib in $(LIBRARIES); do \
		rm -rf "$$lib/reports"; \
		echo "Reportes eliminados: $$lib/reports"; \
	done
