# Project structure

myapp/
├── cmd/                   # entry points for different binaries
│   └── web/               # `go run ./cmd/web`
│       ├── main.go
│       └── config.yaml    # environment‑specific settings
├── internal/              # private application code
│   ├── server/            # HTTP server setup, middleware
│   │   └── server.go
│   ├── routes/            # route definitions
│   │   └── routes.go
│   ├── handlers/          # HTTP handlers
│   │   └── user.go
│   ├── models/            # database models / ORM
│   │   └── user.go
│   ├── migrations/        # SQL or migration scripts
│   └── assets/            # compiled CSS/JS, images
├── templates/             # Templ `.tmpl` files
│   ├── layout.tmpl
│   └── pages/
│       └── index.tmpl
├── web/                   # front‑end build pipeline
│   ├── tailwind.config.js
│   ├── postcss.config.js
│   ├── package.json
│   └── src/               # your CSS/JS entrypoints
├── Dockerfile             # “production” container (Alpine base)
├── docker-compose.yml     # dev environment (optional)
├── Taskfile.yml           # project automation (build, test, watch)
├── Makefile               # optional—for those who like GNU Make
├── go.mod
└── README.md

