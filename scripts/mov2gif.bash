#!/usr/bin/env bash

# Convert mov to an optimized GIF
# Usage: ./mov2gif.bash input.mov [output.gif]

set -e

if [[ $# -lt 1 ]]; then
    echo "Usage: $0 <input.mov> [output.gif]"
    exit 1
fi

input_file="$1"
output_file="${2:-${input_file%.*}.gif}"

if [[ ! -f "$input_file" ]]; then
    echo "Error: File '$input_file' not found"
    exit 1
fi

# Generate color palette
palette="/tmp/palette.png"
options="fps=10,scale=760:-1:flags=lanczos,palettegen=stats_mode=diff"
ffmpeg -i "$input_file" -vf $options -update 1 -y "$palette"

# Convert to GIF using palette
options="fps=10,scale=760:-1:flags=lanczos[x];[x][1:v]paletteuse=dither=sierra2_4a"
ffmpeg -i "$input_file" -i "$palette" -filter_complex "$options" -y "$output_file"

# Cleanup
rm -f "$palette"

echo "Generated: $output_file"