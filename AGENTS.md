# Specialist Agents

Agents pertinents pour le projet view.

## Primary

| Agent | Expertise | Usage |
|-------|-----------|-------|
| `developer-specialist-go` | Go idiomatique, concurrence, error handling | Tout le code applicatif |
| `os-specialist-macos` | macOS, launchd, APFS, SIP, screencapture | Capture d'ecran, permissions |
| `os-specialist-windows-desktop` | Windows 11, winget, firewall | Client Windows, build |

## Supporting

| Agent | Task | Usage |
|-------|------|-------|
| `developer-specialist-review` | Code review (5 sub-executors) | `/review` |
| `developer-executor-correctness` | Invariants, concurrence | Protocole reseau, handshake |
| `developer-executor-security` | Taint analysis, OWASP | Validation du protocole custom |
| `developer-executor-shell` | Shell, Makefile | Build scripts, Makefile |

## Usage

- **Code Go** : `developer-specialist-go` pour tout nouveau code
- **Capture macOS** : `os-specialist-macos` pour les APIs screencapture/CoreGraphics
- **Client Windows** : `os-specialist-windows-desktop` pour les specificitees Windows
- **Review** : `developer-specialist-review` via `/review` avant chaque PR
