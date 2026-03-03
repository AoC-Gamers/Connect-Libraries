#!/usr/bin/env bash
set -euo pipefail

reports_dir="${1:-reports}"
golangci_lint_version="${2:-2.10}"

echo "Ejecutando golangci-lint..."
command -v golangci-lint >/dev/null 2>&1 || {
  echo "Error: golangci-lint no esta instalado. Instalalo con la version ${golangci_lint_version}."
  exit 1
}

installed="$(golangci-lint version 2>/dev/null | head -n1)"
required="${golangci_lint_version#v}"
installed_version="$(echo "$installed" | sed -nE 's/.*version[[:space:]]+v?([0-9]+(\.[0-9]+){1,2}).*/\1/p' | head -n1)"

[ -n "$installed_version" ] || {
  echo "Error: no se pudo detectar la version instalada de golangci-lint."
  exit 1
}

case "$installed_version" in
  "$required"|"$required".*) ;;
  *)
    echo "Error: se requiere golangci-lint ${required}.x (instalada ${installed_version})."
    exit 1
    ;;
esac

required_go="$(sed -nE 's/^go[[:space:]]+([0-9]+(\.[0-9]+){1,2}).*/\1/p' go.mod | head -n1)"
[ -n "$required_go" ] || {
  echo "Error: no se pudo leer la version de Go requerida desde go.mod."
  exit 1
}

echo "$installed" | grep -q "built with go${required_go}" || {
  echo "Error: se requiere golangci-lint compilado con go${required_go}."
  exit 1
}

mkdir -p "$reports_dir"
if golangci-lint run --help 2>&1 | grep -q -- "--output.json.path"; then
  golangci-lint run --timeout=5m --output.json.path "${reports_dir}/lint.json"
else
  golangci-lint run --timeout=5m --out-format json > "${reports_dir}/lint.json"
fi

echo "Reporte JSON generado en ${reports_dir}/lint.json"
