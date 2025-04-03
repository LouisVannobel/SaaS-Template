package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run reset_password.go <email> <new_password>")
		os.Exit(1)
	}

	email := os.Args[1]
	newPassword := os.Args[2]

	// Utiliser des valeurs codées en dur pour la démo
	dbUser := "postgres"
	dbName := "saas_db"
	dbHost := "localhost"
	dbPort := "5432"

	// Construire la chaîne de connexion - pas de mot de passe pour la connexion locale
	connStr := fmt.Sprintf("user=%s dbname=%s host=%s port=%s sslmode=disable", 
		dbUser, dbName, dbHost, dbPort)

	fmt.Println("Tentative de connexion à la base de données...")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Erreur de connexion à la base de données: %v", err)
	}
	defer db.Close()

	// Vérifier la connexion
	err = db.Ping()
	if err != nil {
		log.Fatalf("Erreur de ping à la base de données: %v", err)
	}

	// Hasher le nouveau mot de passe
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Erreur lors du hashage du mot de passe: %v", err)
	}

	// Mettre à jour le mot de passe dans la base de données
	result, err := db.Exec("UPDATE users SET password = $1 WHERE email = $2", string(hashedPassword), email)
	if err != nil {
		log.Fatalf("Erreur lors de la mise à jour du mot de passe: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("Erreur lors de la récupération des lignes affectées: %v", err)
	}

	if rowsAffected == 0 {
		fmt.Printf("Aucun utilisateur trouvé avec l'email: %s\n", email)
	} else {
		fmt.Printf("Mot de passe mis à jour avec succès pour l'utilisateur: %s\n", email)
	}
}
