#!/bin/bash

find . -name "*.go" -o -name "*.js" -o -name "*.html" | entr -r -s "bash scripts/build_styles.sh && go run ./cmd/shortly --config=.application.dev.json"
