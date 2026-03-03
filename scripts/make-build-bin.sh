#!/usr/bin/env bash
set -euo pipefail

output_dir="${1:-bin}"

echo "Compilando binarios en ./${output_dir} ..."
mkdir -p "$output_dir"

for d in cmd/*; do
  if [ -d "$d" ]; then
    name="$(basename "$d")"
    echo "  -> ${name}"
    go build -o "${output_dir}/${name}" "./${d}"
  fi
done

echo "Binarios generados en ./${output_dir}"
