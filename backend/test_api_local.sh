#!/bin/bash

# Variables
API_URL="http://localhost:8081"
AUTH_TOKEN=""
USER_ID=""

# Couleurs pour l'affichage
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "=== Test d'intu00e9gration de l'API SaaS Template (Local) ==="

# Test 1: Santu00e9 de l'API
echo -e "\n1. Test de santu00e9 de l'API"
RESPONSE=$(curl -s -w "\n%{http_code}" $API_URL/health)
HTTP_STATUS=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')

if [ "$HTTP_STATUS" -eq 200 ]; then
    echo -e "${GREEN}u2713 API en ligne: $BODY${NC}"
else
    echo -e "${RED}u2717 API hors ligne: $HTTP_STATUS${NC}"
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
    echo -e "${GREEN}u2713 Inscription ru00e9ussie${NC}"
    USER_ID=$(echo $BODY | grep -o '"id":"[^"]*' | cut -d'"' -f4)
    echo "ID utilisateur: $USER_ID"
else
    echo -e "${RED}u2717 u00c9chec de l'inscription: $HTTP_STATUS${NC}"
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
    echo -e "${GREEN}u2713 Connexion ru00e9ussie${NC}"
    AUTH_TOKEN=$(echo $BODY | grep -o '"token":"[^"]*' | cut -d'"' -f4)
    echo "Token: ${AUTH_TOKEN:0:20}..."
else
    echo -e "${RED}u2717 u00c9chec de la connexion: $HTTP_STATUS${NC}"
    echo "$BODY"
    exit 1
fi

# Test 4: Cru00e9ation d'une tu00e2che
echo -e "\n4. Test de cru00e9ation d'une tu00e2che"
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $AUTH_TOKEN" \
    -d '{"title":"Test Task","description":"This is a test task","status":"pending"}' \
    $API_URL/api/tasks)
HTTP_STATUS=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')

if [ "$HTTP_STATUS" -eq 201 ] || [ "$HTTP_STATUS" -eq 200 ]; then
    echo -e "${GREEN}u2713 Cru00e9ation de tu00e2che ru00e9ussie${NC}"
    TASK_ID=$(echo $BODY | grep -o '"id":"[^"]*' | cut -d'"' -f4)
    echo "ID de la tu00e2che: $TASK_ID"
else
    echo -e "${RED}u2717 u00c9chec de la cru00e9ation de tu00e2che: $HTTP_STATUS${NC}"
    echo "$BODY"
fi

# Test 5: Ru00e9cupu00e9ration des tu00e2ches
echo -e "\n5. Test de ru00e9cupu00e9ration des tu00e2ches"
RESPONSE=$(curl -s -w "\n%{http_code}" -X GET \
    -H "Authorization: Bearer $AUTH_TOKEN" \
    $API_URL/api/tasks)
HTTP_STATUS=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')

if [ "$HTTP_STATUS" -eq 200 ]; then
    echo -e "${GREEN}u2713 Ru00e9cupu00e9ration des tu00e2ches ru00e9ussie${NC}"
    TASKS_COUNT=$(echo $BODY | grep -o '"id"' | wc -l)
    echo "Nombre de tu00e2ches: $TASKS_COUNT"
else
    echo -e "${RED}u2717 u00c9chec de la ru00e9cupu00e9ration des tu00e2ches: $HTTP_STATUS${NC}"
    echo "$BODY"
fi

# Test 6: Profil utilisateur
echo -e "\n6. Test de ru00e9cupu00e9ration du profil utilisateur"
RESPONSE=$(curl -s -w "\n%{http_code}" -X GET \
    -H "Authorization: Bearer $AUTH_TOKEN" \
    $API_URL/api/users/profile)
HTTP_STATUS=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')

if [ "$HTTP_STATUS" -eq 200 ]; then
    echo -e "${GREEN}u2713 Ru00e9cupu00e9ration du profil ru00e9ussie${NC}"
    echo "Email: $(echo $BODY | grep -o '"email":"[^"]*' | cut -d'"' -f4)"
else
    echo -e "${RED}u2717 u00c9chec de la ru00e9cupu00e9ration du profil: $HTTP_STATUS${NC}"
    echo "$BODY"
fi

echo -e "\n=== Tests terminu00e9s ==="
