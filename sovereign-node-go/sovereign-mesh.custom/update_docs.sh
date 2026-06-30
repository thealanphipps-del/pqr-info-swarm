#!/usr/bin/env bash
set -euo pipefail
# Directory containing markdown docs
DOCS_DIR="$(dirname "$0")/docs"
# Add front-matter and TOC to each .md file
for file in $(find "$DOCS_DIR" -type f -name "*.md"); do
  filename=$(basename "$file")
  name_without_ext="${filename%.md}"
  # Generate front-matter if not present
  if ! head -n 1 "$file" | grep -q "^---"; then
    tmp=$(mktemp)
    cat <<EOF > "$tmp"
---
title: "$name_without_ext"
description: ""
date: "2026-05-21"
tags: []
---

EOF
    cat "$file" >> "$tmp"
    mv "$tmp" "$file"
  fi
  # Insert TOC after front-matter (first blank line after ---)
  # Simple TOC: list headings using grep
  toc=$(grep -E "^#{1,6} " "$file" | sed -E "s/^(#{1,6}) (.*)/- [\2](#$(echo \2 | tr '[:upper:]' '[:lower:]' | tr -s ' ' '-' | tr -cd '[:alnum:]-'))/")
  if [ -n "$toc" ]; then
    # Insert TOC after front-matter block
    awk 'BEGIN{found=0} /^---$/ {print; if (found==0) {found=1; next}} {if (found==1 && NF==0) {print "\n## Table of Contents"; print toc; print "\n"; found=2} print}' toc="$toc" "$file" > "${file}.new" && mv "${file}.new" "$file"
  fi
done
echo "Documentation update completed."
