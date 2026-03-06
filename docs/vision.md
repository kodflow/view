# Vision: view

## Purpose

view est un outil de capture d'ecran a distance pour reseau local. Un serveur Go tourne sur macOS, capture l'ecran a la demande, et transmet les images en haute resolution a un client Go multiplateforme qui les affiche via une interface web embarquee.

## Problem Statement

Acceder visuellement a l'ecran de son Mac depuis un PC Windows sur le meme reseau local, sans installer de logiciel lourd (VNC, TeamViewer), sans compte cloud, et sans exposer de service detectable par les scanners reseau.

## Target Users

- Usage personnel uniquement : 1 utilisateur, 2 machines (macOS + Windows)
- Aucun multi-utilisateur prevu

## Goals

1. Capturer l'ecran macOS en haute resolution (4K) sans notification visible
2. Transmettre les captures au client via reseau local avec latence minimale
3. Rester invisible aux scanners reseau (port standard, banniere factice)
4. Client cross-platform testable sur macOS avant deploiement Windows
5. Zero configuration : auto-discovery du serveur sur le reseau local

## Success Criteria

| Critere | Cible |
|---------|-------|
| Resolution capture | Native (Retina/4K) |
| Latence capture-affichage | < 500ms en LAN |
| Detection par nmap -sV | Identifie comme SSH, pas comme service custom |
| Binaires | 2 binaires statiques (server macOS, client multi-OS) |
| Setup | Zero config, auto-discovery mDNS/UDP broadcast |
| Interface client | Web UI embarquee, ouvre le navigateur |

## Design Principles

- **Simplicite** : un binaire serveur, un binaire client, zero dependance externe
- **Discretion reseau** : port 22 (ou standard), banniere SSH factice pour les connexions inconnues, magic bytes handshake pour le client legitime
- **Cross-platform** : le client compile et fonctionne sur macOS et Windows sans modification
- **Confiance locale** : pas d'authentification, le reseau local est le perimetre de securite
- **Capture native** : utiliser les APIs macOS (screencapture CLI ou CoreGraphics) pour la meilleure qualite

## Non-Goals

- Streaming video temps reel (pseudo-streaming par captures rapides suffit)
- Support multi-utilisateurs ou multi-serveurs
- Authentification par mot de passe ou certificat
- Fonctionnement hors reseau local (pas d'Internet/WAN)
- Controle a distance (clavier/souris) - lecture seule uniquement
- Interface GUI native (web UI embarquee suffit)

## Key Decisions

| Decision | Choix | Raison |
|----------|-------|--------|
| Langage | Go | Expertise existante, cross-compile trivial, binaires statiques |
| Capture ecran | screencapture CLI / CoreGraphics | APIs macOS natives, pas de notification |
| Transport | TCP sur port standard (22) | Invisible aux scanners, banniere SSH factice |
| Discovery | mDNS ou UDP broadcast | Zero config, reseau local uniquement |
| Interface client | Web UI embarquee (HTML/JS) | Cross-platform, pas de dep GUI |
| Format image | PNG ou WebP | Qualite maximale, compression efficace |
| Handshake | Magic bytes custom | Distingue client legitime des scanners |
