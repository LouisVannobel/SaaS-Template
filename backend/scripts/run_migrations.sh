#!/bin/bash

# Récupérer les variables d'environnement
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-postgres_password}
DB_NAME=${DB_NAME:-saas_db}

echo "Exécution des migrations SQL..."

# Exécuter le script de migration
PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f ./db/migrations/000001_init_schema.up.sql

if [ $? -eq 0 ]; then
    echo "Migrations exécutées avec succès!"
else
    echo "Erreur lors de l'exécution des migrations."
    exit 1
fi
