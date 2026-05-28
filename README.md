# shcc - Share Claude Code

**shcc** est un outil CLI et une API permettant de partager de manière sécurisée vos credentials [Claude Code](https://docs.anthropic.com/claude/docs/claude-code) entre différentes machines (ex: de votre ordinateur personnel vers une VM cloud).

## 🚀 Fonctionnalités

- **Partage sécurisé** : Utilise un chiffrement hybride (RSA-4096 + AES-256-GCM). Vos credentials ne sont jamais lisibles en clair par le serveur.
- **Identification simplifiée** : Partagez via email ou via un pseudonyme personnalisé (`shcc name`).
- **Multi-OS** : Support complet de Linux, macOS et Windows (WSL et Natif).
- **Synchronisation automatique** : Votre clé publique et votre profil sont mis à jour à chaque utilisation.

## 🏗️ Architecture (Monorepo)

Le projet est structuré en Go pour séparer le client de l'API tout en partageant la logique de sécurité :

- `cmd/shcc-cli/` : Code source du CLI.
- `cmd/shcc-api/` : Service API HTTP conçu pour Google Cloud Run.
- `internal/auth/` : Cœur de la cryptographie (RSA/AES).
- `internal/database/` : Couche d'accès aux données PostgreSQL.
- `internal/config/` : Gestion des configurations Claude Code et shcc.

## 🛠️ Installation du CLI

### Via script (recommandé)
```bash
curl -fsSL https://shcc.ongoua.pro/install.sh | bash
```

### Manuellement (depuis les sources)
```bash
go build -o shcc ./cmd/shcc-cli
mv shcc /usr/local/bin/
```

## 📖 Utilisation

### 1. Vérifier votre statut
```bash
shcc status
```

### 2. Définir un pseudonyme
```bash
shcc name didier
```

### 3. Partager ses credentials
```bash
shcc didier.vincent@gmail.com
# ou si Didier a configuré son nom :
shcc didier
```

### 4. Recevoir et installer des credentials
Lancez simplement `shcc` sans argument sur la machine de destination :
```bash
shcc
```
Le CLI listera les clés reçues et vous proposera de les installer.

## ☁️ Déploiement de l'API (Cloud Run)

L'API nécessite une base de données PostgreSQL existante.

### 1. Appliquer le schéma
Exécutez le fichier `internal/database/schema.sql` sur votre instance PostgreSQL.

### 2. Déployer sur Google Cloud Run
```bash
# Build de l'image
gcloud builds submit --tag gcr.io/[PROJECT_ID]/shcc-api ./cmd/shcc-api

# Déploiement
gcloud run deploy shcc-api \
  --image gcr.io/[PROJECT_ID]/shcc-api \
  --platform managed \
  --region europe-west1 \
  --set-env-vars "DATABASE_URL=postgres://user:pass@host:5432/db" \
  --allow-unauthenticated
```

## 🔒 Sécurité

Le processus de partage garantit une confidentialité totale :
1. Chaque client génère une paire de clés RSA-4096 localement (`~/.shcc/id_rsa`).
2. La clé publique est envoyée au serveur.
3. Lors d'un partage, l'expéditeur récupère la clé publique du destinataire.
4. Les credentials sont chiffrés avec AES-256, et la clé AES est chiffrée avec la clé publique RSA du destinataire.
5. Seul le destinataire possède la clé privée capable de déchiffrer le message final.

## 📄 Licence
MIT - Créé par Dimitri Ongoua
