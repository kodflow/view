# Architecture: view

## System Context

```
┌─────────────────────┐        LAN (TCP:22)       ┌─────────────────────┐
│   macOS Server      │◄─────────────────────────►│   Client (Win/Mac)  │
│                     │   magic bytes handshake    │                     │
│  ┌───────────────┐  │   PNG/WebP frames          │  ┌───────────────┐  │
│  │ Screen Capture│  │                            │  │ Web UI        │  │
│  │ (native API)  │  │                            │  │ (embed.FS)    │  │
│  └───────────────┘  │                            │  └───────────────┘  │
│                     │   mDNS / UDP broadcast     │                     │
│  ┌───────────────┐  │◄─────────────────────────►│  ┌───────────────┐  │
│  │ Stealth Port  │  │   auto-discovery           │  │ Discovery     │  │
│  │ (SSH banner)  │  │                            │  │ (listener)    │  │
│  └───────────────┘  │                            │  └───────────────┘  │
└─────────────────────┘                            └─────────────────────┘
                                                     │
                                                     ▼
                                                   Browser
                                                   (localhost)
```

## Components

### Server (macOS)

| Module | Responsabilite |
|--------|----------------|
| `capture` | Capture d'ecran via screencapture CLI ou CoreGraphics CGo |
| `network` | Ecoute TCP sur port 22, banniere SSH factice, handshake magic bytes |
| `discovery` | Annonce mDNS/UDP broadcast pour auto-detection par le client |

### Client (cross-platform)

| Module | Responsabilite |
|--------|----------------|
| `discovery` | Detecte le serveur via mDNS/UDP broadcast sur le LAN |
| `network` | Connexion TCP, envoi magic bytes, reception des frames |
| `web` | Serveur HTTP local (localhost), sert la Web UI embarquee |

### Shared

| Module | Responsabilite |
|--------|----------------|
| `protocol` | Definition du protocole : magic bytes, headers, commandes |

## Data Flow

1. **Discovery** : client envoie mDNS query ou UDP broadcast → serveur repond avec son IP
2. **Connect** : client ouvre TCP vers serveur:22 → envoie magic bytes
3. **Handshake** : serveur verifie magic bytes → bascule en mode capture (sinon banniere SSH)
4. **Capture** : client envoie commande "capture" → serveur capture l'ecran → encode PNG/WebP → envoie la frame
5. **Display** : client recoit la frame → la sert via HTTP local → le navigateur l'affiche

## Technology Stack

| Couche | Technologie |
|--------|-------------|
| Langage | Go 1.25+ |
| Capture ecran | screencapture (CLI) ou CoreGraphics (CGo) |
| Transport | TCP raw, protocole binaire custom |
| Discovery | mDNS (Bonjour) ou UDP broadcast |
| Encodage image | PNG (qualite max) ou WebP (compression) |
| Client UI | HTML/JS embarque via Go embed.FS |
| Build | Makefile, GOOS/GOARCH cross-compile |

## Constraints

- Le serveur ne tourne que sur macOS (APIs de capture specifiques)
- Le client doit fonctionner sur macOS ET Windows sans recompilation conditionnelle
- Pas de dependance CGo cote client (pur Go pour la portabilite)
- Le port 22 necessite des privileges root/admin sur macOS (ou port > 1024 en fallback)
