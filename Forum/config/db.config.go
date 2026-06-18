package config

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// InitDB ouvre la base MySQL et vérifie qu'elle répond bien.
func InitDB() *sql.DB {
	dsn := GetEnv("DB_DSN", "root:password@tcp(127.0.0.1:3306)/forum?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci&loc=Local")

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Erreur ouverture base de données : %s", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Erreur connexion base de données : %s", err)
	}

	log.Println("Base de données connectée via MySQL")
	return db
}
