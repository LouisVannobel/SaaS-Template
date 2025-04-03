#!/bin/bash

# Variables
API_URL="http://localhost:8080"
AUTH_TOKEN=""
USER_ID=""

# Couleurs pour l'affichage
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "=== Test d'intégration de l'API SaaS Template ==="

# Test 1: Santé de l'API
echo -e "\n1. Test de santé de l'API"
RESPONSE=$(curl -s -w "\n%{http_code}" $API_URL/health)
HTTP_STATUS=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')

if [ "$HTTP_STATUS" -eq 200 ]; then
    echo -e "${GREEN}✓ API en ligne: $BODY${NC}"
else
    echo -e "${RED}✗ API hors ligne: $HTTP_STATUS${NC}"
    exit 1
fi

# Test 2: Inscription d'un utilisateur
echo -e "\n2. Test d'inscription d'un utilisateur"
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST \
    -H "Content-Type: application/json" \
    -d '{"email":"test@example.com","password":"Password123!","name":"Test User"}' \
    $API_URL/api/register)
HTTP_STATUS=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')

if [ "$HTTP_STATUS" -eq 201 ] || [ "$HTTP_STATUS" -eq 200 ]; then
    echo -e "${GREEN}✓ Inscription réussie${NC}"
    USER_ID=$(echo $BODY | grep -o '"id":"[^"]*' | cut -d'"' -f4)
    echo "ID utilisateur: $USER_ID"
else
    echo -e "${RED}✗ Échec de l'inscription: $HTTP_STATUS${NC}"
    echo "$BODY"
fi

# Test 3: Connexion de l'utilisateur
echo -e "\n3. Test de connexion"
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST \
    -H "Content-Type: application/json" \
    -d '{"email":"test@example.com","password":"Password123!"}' \
    $API_URL/api/login)
HTTP_STATUS=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')

if [ "$HTTP_STATUS" -eq 200 ]; then
    echo -e "${GREEN}✓ Connexion réussie${NC}"
    AUTH_TOKEN=$(echo $BODY | grep -o '"token":"[^"]*' | cut -d'"' -f4)
    echo "Token: ${AUTH_TOKEN:0:20}..."
else
    echo -e "${RED}✗ Échec de la connexion: $HTTP_STATUS${NC}"
    echo "$BODY"
    exit 1
fi

# Test 4: Création d'une tâche
echo -e "\n4. Test de création d'une tâche"
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $AUTH_TOKEN" \
    -d '{"title":"Test Task","description":"This is a test task","status":"pending"}' \
    $API_URL/api/tasks)
HTTP_STATUS=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')

if [ "$HTTP_STATUS" -eq 201 ] || [ "$HTTP_STATUS" -eq 200 ]; then
    echo -e "${GREEN}✓ Création de tâche réussie${NC}"
    TASK_ID=$(echo $BODY | grep -o '"id":"[^"]*' | cut -d'"' -f4)
    echo "ID de la tâche: $TASK_ID"
else
    echo -e "${RED}✗ Échec de la création de tâche: $HTTP_STATUS${NC}"
    echo "$BODY"
fi

# Test 5: Récupération des tâches
echo -e "\n5. Test de récupération des tâches"
RESPONSE=$(curl -s -w "\n%{http_code}" -X GET \
    -H "Authorization: Bearer $AUTH_TOKEN" \
    $API_URL/api/tasks)
HTTP_STATUS=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')

if [ "$HTTP_STATUS" -eq 200 ]; then
    echo -e "${GREEN}✓ Récupération des tâches réussie${NC}"
    TASKS_COUNT=$(echo $BODY | grep -o '"id"' | wc -l)
    echo "Nombre de tâches: $TASKS_COUNT"
else
    echo -e "${RED}✗ Échec de la récupération des tâches: $HTTP_STATUS${NC}"
    echo "$BODY"
fi

# Test 6: Profil utilisateur
echo -e "\n6. Test de récupération du profil utilisateur"
RESPONSE=$(curl -s -w "\n%{http_code}" -X GET \
    -H "Authorization: Bearer $AUTH_TOKEN" \
    $API_URL/api/users/profile)
HTTP_STATUS=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')

if [ "$HTTP_STATUS" -eq 200 ]; then
    echo -e "${GREEN}✓ Récupération du profil réussie${NC}"
    echo "Email: $(echo $BODY | grep -o '"email":"[^"]*' | cut -d'"' -f4)"
else
    echo -e "${RED}✗ Échec de la récupération du profil: $HTTP_STATUS${NC}"
    echo "$BODY"
fi

echo -e "\n=== Tests terminés ==="
