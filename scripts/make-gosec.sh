#!/usr/bin/env bash
set -euo pipefail

reports_dir="${1:-reports}"
gosec_version="${2:-2.23}"

echo "Ejecutando escaneo de seguridad con gosec..."
command -v gosec >/dev/null 2>&1 || {
  echo "Error: gosec no esta instalado. Instalalo con la version ${gosec_version}."
  exit 1
}

mkdir -p "$reports_dir"
rc=0
gosec -fmt=json -out="${reports_dir}/gosec.json" ./... || rc=$?

echo "Reporte JSON generado en ${reports_dir}/gosec.json"

if command -v jq >/dev/null 2>&1; then
  found="$(jq -r '.Stats.found // 0' "${reports_dir}/gosec.json")"
  files="$(jq -r '.Stats.files // 0' "${reports_dir}/gosec.json")"
  nosec="$(jq -r '.Stats.nosec // 0' "${reports_dir}/gosec.json")"
  printf "Resumen gosec: found=%s files=%s nosec=%s\n" "$found" "$files" "$nosec"
else
  found="$(grep -o '"rule_id"' "${reports_dir}/gosec.json" | wc -l | tr -d '[:space:]')"
  printf "Resumen gosec: found=%s (instala jq para ver mas detalle)\n" "$found"
fi

if [ "$rc" -ne 0 ]; then
  echo "gosec finalizo con errores (exit ${rc}). Revisa ${reports_dir}/gosec.json"
  exit "$rc"
fi
