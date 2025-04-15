package engine

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
    var err error

    dbUser := os.Getenv("DB_USER")
    dbPass := os.Getenv("DB_PASS")
    dbName := os.Getenv("DB_NAME")
    dbHost := os.Getenv("DB_HOST")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbUser, dbPass, dbHost, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Erreur de connexion à MySQL:", err)
	}
	defer db.Close()

    DB, err = sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal("Erreur de connexion à MySQL :", err)
    }   

	fmt.Println("Connexion réussie à MySQL 🎉")
}