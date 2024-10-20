#!/usr/bin/env zsh

echo "Building styles..."
./tailwindcss -i ./cmd/shortly/main.css -o ./cmd/shortly/public/css/output.css --minify
echo "Styles built!"