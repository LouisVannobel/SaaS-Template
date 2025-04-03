#!/bin/bash

echo "Exu00e9cution des migrations SQL via Docker..."

# Copier le fichier de migration dans le conteneur PostgreSQL
sudo docker cp ./db/migrations/000001_init_schema.up.sql saas_postgres:/tmp/

# Exu00e9cuter le script de migration dans le conteneur
sudo docker exec saas_postgres psql -U postgres -d saas_db -f /tmp/000001_init_schema.up.sql

if [ $? -eq 0 ]; then
    echo "Migrations exu00e9cutu00e9es avec succu00e8s!"
else
    echo "Erreur lors de l'exu00e9cution des migrations."
    exit 1
fi
