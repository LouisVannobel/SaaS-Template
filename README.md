# SaaS Template

Un template d'application SaaS multi-tiers avec une architecture moderne et évolutive.

## Aperçu du Projet

Ce projet est un template complet pour le développement d'applications SaaS (Software as a Service) avec une architecture multi-tiers comprenant :

- **Base de données** : PostgreSQL pour le stockage persistant des données
- **Backend** : API REST en Go avec authentification JWT
- **Frontend** : Interface utilisateur réactive en React avec TypeScript

L'ensemble est conteneurisé avec Docker pour faciliter le déploiement et le développement.

## Fonctionnalités

### Authentification et Gestion des Utilisateurs
- Inscription et connexion sécurisées
- Authentification basée sur JWT
- Hachage des mots de passe en SHA256
- Profil utilisateur personnalisable

### Gestion des Tâches
- Création, lecture, mise à jour et suppression de tâches
- Filtrage et tri des tâches
- Attribution de tâches à des utilisateurs

### Architecture Technique
- API RESTful avec documentation
- Validation des données entrantes
- Gestion des erreurs robuste
- Journalisation complète (console + fichiers)

## Architecture

### Base de Données (PostgreSQL)
- Tables pour les utilisateurs et les tâches
- Relations et contraintes pour l'intégrité des données
- Indexation pour des performances optimales

### Backend (Go)
- Structure de projet organisée et modulaire
- Utilisation de goroutines pour les opérations asynchrones
- Channels pour la communication entre composants
- Middleware d'authentification et de validation

### Frontend (React + TypeScript)
- Architecture basée sur les composants
- Gestion d'état avec Context API
- Routes protégées pour les utilisateurs authentifiés
- Interface utilisateur responsive avec Tailwind CSS

## Installation et Démarrage

### Prérequis
- Docker et Docker Compose
- Git

### Installation

1. Cloner le dépôt :
   ```bash
   git clone https://github.com/LouisVannobel/SaaS-Template.git
   cd SaaS-Template
   ```

2. Créer un fichier `.env` à partir du modèle `.env.example` :
   ```bash
   cp .env.example .env
   ```

3. Modifier les variables d'environnement selon vos besoins dans le fichier `.env`

### Démarrage avec Docker

1. Construire et démarrer les conteneurs :
   ```bash
   docker-compose up -d
   ```

2. Exécuter les migrations de base de données :
   ```bash
   cd backend && ./scripts/run_migrations_docker.sh
   ```

3. Accéder à l'application :
   - Frontend : http://localhost:FRONTEND_PORT
   - API Backend : http://localhost:BACKEND_PORT
   - Base de données : localhost:DB_PORT
   

### Développement Local

#### Backend
```bash
cd backend
export DB_HOST=localhost
go run main.go
```

#### Frontend
```bash
cd frontend
npm install
npm run dev
```

## Structure du Projet

```
├── backend/                # Code source du backend Go
│   ├── api/                # Handlers, routes et middleware API
│   ├── auth/               # Authentification et sécurité
│   ├── db/                 # Couche d'accès aux données et migrations
│   ├── models/             # Modèles de données
│   └── utils/              # Utilitaires (logging, etc.)
├── frontend/              # Code source du frontend React
│   ├── public/             # Ressources statiques
│   └── src/                # Code source React
│       ├── components/     # Composants React réutilisables
│       ├── contexts/       # Contextes React pour la gestion d'état
│       ├── pages/          # Pages de l'application
│       └── services/       # Services pour les appels API
├── .env.example           # Modèle de variables d'environnement
├── docker-compose.yml     # Configuration Docker Compose
├── backend.Dockerfile     # Dockerfile pour le backend
└── frontend.Dockerfile    # Dockerfile pour le frontend
```

## Déploiement

Pour un déploiement en production, consultez le fichier [DEPLOYMENT.md](DEPLOYMENT.md) qui contient des instructions détaillées.

## Licence

Ce projet est sous licence GNU AFFERO GENERAL PUBLIC LICENSE. Voir le fichier [LICENSE](LICENSE) pour plus de détails.
