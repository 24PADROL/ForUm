package engine

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql" // Remplace par ton driver si nécessaire
	"golang.org/x/crypto/bcrypt"
)

// Page d'accueil
func HomeHandler(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("web/html/home.html"))
    tmpl.Execute(w, nil)
}

// Inscription
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        tmpl, err := template.ParseFiles("/web/html/register.html")
        if err != nil {
            http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
            return
        }
        tmpl.Execute(w, nil)
        return
    }

	if r.Method == "POST" {
		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			log.Println("Erreur de décodage JSON:", err)
			http.Error(w, "Données invalides", http.StatusBadRequest)
			return
		}

		// Hachage du mot de passe avec gestion d'erreur
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Erreur lors du hachage du mot de passe:", err)
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}

        query := `INSERT INTO users (username, email, password) VALUES (?, ?, ?)`
        _, err = DB.Exec(query, user.Username, user.Email, hashedPassword)
        if err != nil {
            log.Println("Erreur lors de l'insertion dans la base de données:", err)
            http.Error(w, "Erreur lors de l'inscription", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(map[string]string{"message": "Inscription réussie"})
    } else {
        http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
    }
}

// Connexion
func LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        tmpl, err := template.ParseFiles("web/html/login.html")
        if err != nil {
            http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
            return
        }
        tmpl.Execute(w, nil)
        return
    }

	if r.Method == http.MethodPost {
		var user User

		// Lire les données du formulaire
		body, _ := io.ReadAll(r.Body)
		log.Println("Corps de la requête reçu:", string(body)) // Log des données du formulaire

		// Décoder les données JSON
		if err := json.Unmarshal(body, &user); err != nil {
			log.Println("Erreur de décodage JSON:", err)
			http.Error(w, "Données invalides", http.StatusBadRequest)
			return
		}

		// Vérification du format des données reçues
		log.Println("Données reçues après décodage:", user)

        // Rediriger vers la page d'accueil
        http.Redirect(w, r, "/accueil", http.StatusSeeOther)
    }
}

func AccueilHandler(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("web/html/accueil.html"))
    tmpl.Execute(w, nil)
}