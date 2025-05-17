#!/bin/bash

# Assumes:
# {.air.toml, docker-compose.yml, .git, .gitignore, go.mod, README.md, Taskfile.yml}

echo "--> Initializing Go web app project structure..."

# Directories
mkdir -p cmd/web
mkdir -p internal/{server,routes,handlers,models,migrations,assets/css}
mkdir -p templates/pages
mkdir -p web/src

# Base files
touch cmd/web/main.go
touch cmd/web/config.yaml
touch internal/server/server.go
touch internal/routes/routes.go
touch internal/handlers/handlers.go
touch internal/models/models.go
touch templates/layout.tmpl
touch templates/pages/index.tmpl
touch web/src/input.css
touch web/tailwind.config.js
touch web/postcss.config.js
touch web/package.json

# Initialize README if missing
if [ ! -f README.md ]; then
	echo "# My Web App" > README.md
fi

# Add sample main.go if empty
if [ ! -s cmd/web/main.go ]; then
	cat <<EOF > cmd/web/main.go
package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
EOF
fi

# Add example Tailwind input file
echo '@tailwind base; @tailwind components; @tailwind utilities;' > web/src/input.css

# Success message
echo "✅ Project scaffold created!"
