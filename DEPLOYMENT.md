# Déploiement du SaaS Template

Ce document fournit des instructions détaillées pour déployer l'application SaaS Template dans différents environnements.

## Déploiement en Production

### Prérequis

- Un serveur Linux (Ubuntu 22.04 LTS recommandé)
- Docker et Docker Compose installés
- Un nom de domaine configuré
- Certificats SSL (Let's Encrypt recommandé)

### Étapes de déploiement

1. **Cloner le dépôt sur le serveur**

   ```bash
   git clone https://github.com/LouisVannobel/SaaS-Template.git
   cd SaaS-Template
   ```

2. **Configurer les variables d'environnement**

   Créez un fichier `.env` basé sur `.env.example` et modifiez les valeurs pour l'environnement de production :

   ```bash
   cp .env.example .env
   nano .env
   ```

   Variables importantes à modifier :
   - `DB_PASSWORD` : Utilisez un mot de passe fort
   - `JWT_SECRET` : Générez une clé secrète unique
   - `FRONTEND_PORT` et `BACKEND_PORT` : Configurez selon vos besoins

3. **Configurer Nginx comme proxy inverse**

   Créez un fichier de configuration Nginx pour votre domaine :

   ```bash
   sudo nano /etc/nginx/sites-available/saas-template.conf
   ```

   Exemple de configuration :

   ```nginx
   server {
       listen 80;
       server_name votre-domaine.com;
       return 301 https://$host$request_uri;
   }

   server {
       listen 443 ssl;
       server_name votre-domaine.com;

       ssl_certificate /chemin/vers/certificat.pem;
       ssl_certificate_key /chemin/vers/cle-privee.pem;

       # Frontend
       location / {
           proxy_pass http://localhost:3000;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
       }

       # Backend API
       location /api {
           proxy_pass http://localhost:8080;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
       }
   }
   ```

   Activez la configuration :

   ```bash
   sudo ln -s /etc/nginx/sites-available/saas-template.conf /etc/nginx/sites-enabled/
   sudo nginx -t
   sudo systemctl reload nginx
   ```

4. **Démarrer les conteneurs Docker**

   ```bash
   docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
   ```

5. **Exécuter les migrations de base de données**

   ```bash
   cd backend && ./scripts/run_migrations_docker.sh
   ```

6. **Vérifier le déploiement**

   Accédez à votre domaine dans un navigateur pour vérifier que l'application fonctionne correctement.

## Déploiement sur Heroku

Pour déployer sur Heroku, suivez les étapes suivantes :

1. **Installer l'interface en ligne de commande Heroku**

   Suivez les instructions sur [devcenter.heroku.com](https://devcenter.heroku.com/articles/heroku-cli)

2. **Se connecter à Heroku**

   ```bash
   heroku login
   ```

3. **Créer une application Heroku**

   ```bash
   heroku create votre-app-saas
   ```

4. **Ajouter une base de données PostgreSQL**

   ```bash
   heroku addons:create heroku-postgresql:hobby-dev
   ```

5. **Configurer les variables d'environnement**

   ```bash
   heroku config:set JWT_SECRET=votre_secret_jwt
   ```

6. **Déployer l'application**

   ```bash
   git push heroku main
   ```

7. **Exécuter les migrations**

   ```bash
   heroku run bash -a votre-app-saas
   cd backend && go run scripts/migrate.go up
   ```

## Mise à jour de l'application

Pour mettre à jour l'application déjà déployée :

1. **Extraire les dernières modifications**

   ```bash
   git pull origin main
   ```

2. **Reconstruire et redémarrer les conteneurs**

   ```bash
   docker-compose down
   docker-compose up -d --build
   ```

3. **Exécuter les nouvelles migrations si nécessaire**

   ```bash
   cd backend && ./scripts/run_migrations_docker.sh
   ```
