#!/bin/bash
set -e

# Ensure templ is available
if ! command -v templ &> /dev/null
then
    if [ -f ~/go/bin/templ ]; then
        TEMPL=~/go/bin/templ
    else
        echo "templ could not be found, please install it with 'go install github.com/a-h/templ/cmd/templ@latest'"
        exit 1
    fi
else
    TEMPL=command -v templ
fi

echo "Generating Templ templates..."
$TEMPL generate

echo "Compiling Go to WASM..."
# Note: Using standard go build as tinygo was not found in the environment
GOOS=js GOARCH=wasm go build -o dist/main.wasm ./wasm/main.go

echo "Compiling Tailwind CSS..."
npx @tailwindcss/cli -i ./styles/input.css -o ./dist/output.css

echo "Copying static files..."
cp static/* dist/
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" dist/

echo "Build complete! Artifacts are in the 'dist' directory."
