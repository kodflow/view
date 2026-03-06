# Development Workflows

## Setup

### Prerequisites

- Go 1.25+ installe
- macOS pour le serveur (capture d'ecran)
- Make

### Installation

```bash
git clone https://github.com/kodflow/view.git
cd view
make deps
```

## Development Loop

```bash
# 1. Editer le code dans src/
# 2. Tester
make test

# 3. Builder le serveur (macOS)
make build-server

# 4. Builder le client (local)
make build-client

# 5. Lancer le serveur
./bin/view-server

# 6. Lancer le client (autre terminal)
./bin/view-client
```

## Testing Strategy

### Unit Tests

- `src/internal/capture/` : mock de screencapture, test d'encodage
- `src/internal/network/` : test du handshake, banniere SSH, protocole
- `src/internal/discovery/` : test mDNS/broadcast avec loopback

### Integration Tests

- Serveur + client en loopback sur localhost
- Test du cycle complet : discovery → connect → capture → display

## Build

### Binaires locaux

```bash
make build          # Build server + client pour l'OS courant
make build-server   # Serveur uniquement
make build-client   # Client uniquement
```

### Cross-compile

```bash
make build-all      # Tous les binaires
# Produit :
#   bin/view-server-darwin-arm64
#   bin/view-client-darwin-arm64
#   bin/view-client-windows-amd64.exe
```

## CI/CD

### Pipeline

1. `make lint` — golangci-lint
2. `make test` — go test avec race detector
3. `make build-all` — cross-compilation
4. Release — binaires attaches au tag GitHub
