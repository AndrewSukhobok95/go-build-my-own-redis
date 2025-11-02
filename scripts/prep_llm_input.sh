#!/usr/bin/env bash
set -euo pipefail

# ==========================
# CONFIGURATION
# ==========================

# Optional first argument: project root folder
# If not provided, default to one directory up from the script
if [ $# -ge 1 ]; then
    PROJECT_ROOT="$1"
else
    PROJECT_ROOT="$(dirname "$(realpath "$0")")/.."
fi

# Resolve PROJECT_ROOT to an absolute path
PROJECT_ROOT="$(realpath "$PROJECT_ROOT")"

echo "PROJECT_ROOT is set to ${PROJECT_ROOT}" 

# Output files (will be created in project root)
ALL_FILES="$PROJECT_ROOT/all_go_files.txt"
FUNC_SIGS="$PROJECT_ROOT/signatures.txt"

# =================================================================
# Step 1: Generate all_go_files.txt
# =================================================================
echo "Generating ${ALL_FILES} from Go files in ${PROJECT_ROOT} ..."
> "$ALL_FILES"
find "$PROJECT_ROOT" -type f -name "*.go" | while read -r file; do
  echo "===== ${file} =====" >> "$ALL_FILES"
  cat "$file" >> "$ALL_FILES"
  echo -e "\n" >> "$ALL_FILES"
done

# =================================================================
# Step 2: Extract function signatures and struct definitions
# =================================================================
echo "Extracting function signatures and structs into ${FUNC_SIGS} ..."
awk '
# Detect file path header lines like: ===== ./path/to/file.go =====
/^===== .* =====$/ {
    print
    next
}

# Match single-line function definitions
/^func[[:space:]]+[A-Za-z0-9_]+\([^)]*\)[[:space:]]*[A-Za-z0-9_.*\[\]]*[[:space:]]*{$/ {
    sub(/[[:space:]]*{$/, "", $0)
    print
    next
}

# Detect start of a type struct block
/^type[[:space:]]+[A-Za-z0-9_]+[[:space:]]+struct[[:space:]]*{$/ {
    in_struct = 1
    struct_def = $0
    next
}

# Capture lines inside the struct until closing brace
in_struct {
    struct_def = struct_def "\n" $0
    if ($0 ~ /^}/) {
        print struct_def
        in_struct = 0
    }
    next
}
' "$ALL_FILES" > "$FUNC_SIGS"

echo "Done!"
echo "Created:"
echo " - $ALL_FILES"
echo " - $FUNC_SIGS"
