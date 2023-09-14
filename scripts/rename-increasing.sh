#!/usr/bin/env bash
# Renames all files in the current directory to contain a chapter number
# whenever the last digits in the filename are a number and reset.
# Starts at 1 and increases by 1 for each file.
set -euo pipefail

title="$(basename "$(pwd)")"

# Loop through all files
chapter=1
last_rel_number=0
for file in *; do
  # Expected format:
  # ABS_REL.EXT
  ext="${file##*.}"
  abs="$(awk -F '_' '{print $1}' <<<"$file")"
  rel_number="${file##*_}"
  rel_number="${rel_number%.*}"

  if [ "$rel_number" -lt "$last_rel_number" ]; then
    chapter=$((chapter + 1))
  fi

  new_file="${abs} $title v1 c${chapter} p$rel_number.${ext}"

  #echo "$file" "->" "$new_file"
  mv -v "$file" "$new_file"

  last_rel_number="$rel_number"
done
