# view

## Purpose

Outil de capture d'ecran macOS a distance via reseau local. Serveur Go sur macOS, client Go multiplateforme avec web UI embarquee. Auto-discovery, port stealth (banniere SSH factice), captures haute resolution.

## Project Structure

```
/workspace
├── src/
│   ├── cmd/
│   │   ├── server/          # Point d'entree serveur (macOS)
│   │   └── client/          # Point d'entree client (cross-platform)
│   ├── internal/
│   │   ├── capture/         # Capture d'ecran (screencapture/CoreGraphics)
│   │   ├── network/         # Transport TCP, handshake, banniere SSH
│   │   ├── discovery/       # mDNS/UDP broadcast auto-discovery
│   │   └── web/             # Serveur web embarque + assets HTML/JS
│   └── go.mod
├── docs/                    # Documentation projet
├── CLAUDE.md                # Ce fichier
└── Makefile                 # Build targets
```

## Tech Stack

- **Language**: Go 1.25+
- **Capture**: screencapture CLI / CoreGraphics (CGo)
- **Transport**: TCP raw sur port 22, protocole custom avec magic bytes
- **Discovery**: mDNS (Bonjour) ou UDP broadcast
- **Client UI**: HTML/JS embarque via embed.FS
- **Build**: Makefile, cross-compile GOOS/GOARCH

## How to Work

1. **New feature**: `/plan "description"` → `/do` → `/git --commit`
2. **Bug fix**: `/plan "description"` → `/do` → `/git --commit`
3. **Test**: `make test` ou `go test ./src/...`
4. **Build**: `make build` (binaires dans `bin/`)
5. **Cross-compile**: `make build-all` (macOS server + Windows/macOS client)

Branch conventions: `feat/<desc>` ou `fix/<desc>`.

## Key Principles

- **Go idiomatique** : error handling explicite, interfaces minimales, pas de frameworks
- **Binaires statiques** : zero dependance runtime
- **MCP-first** : utiliser les outils MCP avant les CLI
- **Semantic search** : grepai pour les recherches, Grep pour les chaines exactes

## Verification

- `make test` : tests unitaires
- `make lint` : golangci-lint
- `make build` : compilation sans erreur
- Pas de secrets dans les commits (hook security)
- Format conventionnel pour les commits
