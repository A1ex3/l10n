#!/usr/bin/env bash
set -euo pipefail

FILE="${MODULE_FILE:-MODULE.bazel}"

usage() {
  echo "Use: $0 -name | -version [path/to/MODULE.bazel]" >&2
  exit 1
}

[[ $# -ge 1 ]] || usage

case "$1" in
  -name|--name)       KEY="name" ;;
  -version|--version) KEY="version" ;;
  *) usage ;;
esac

[[ $# -ge 2 ]] && FILE="$2"
[[ -f "$FILE" ]] || { echo "Файл не найден: $FILE" >&2; exit 1; }

BLOCK=$(awk '
  /^[[:space:]]*module[[:space:]]*\(/ { f = 1 }
  f { print }
  f && /\)/ { exit }
' "$FILE")

VALUE=$(printf '%s\n' "$BLOCK" \
  | sed -nE "s/.*(^|[^[:alnum:]_])${KEY}[[:space:]]*=[[:space:]]*\"([^\"]*)\".*/\2/p" \
  | head -n1)

if [[ -z "$VALUE" ]]; then
  echo "Field '$KEY' not found at module() file: $FILE" >&2
  exit 1
fi

echo "$VALUE"