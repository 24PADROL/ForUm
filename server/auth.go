package engine

import (
	"database/sql"
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
	tmpl := template.Must(template.ParseFiles("template/Home.tmpl"))
	tmpl.Execute(w, nil)
}

// Inscription
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("template/register.tmpl")
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

		// Insertion dans la base de données
		query := `INSERT INTO users (username, email, password) VALUES (?, ?, ?)`
		_, err = DB.Exec(query, user.Username, user.Email, hashedPassword)
		if err != nil {
			log.Println("Erreur lors de l'insertion utilisateur:", err)
			http.Error(w, "Erreur lors de l'inscription", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Inscription réussie"})
	}
}

// Connexion
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("template/login.html")
		if err != nil {
			log.Println("Erreur de parsing du template:", err) // Log d'erreur pour le template
			http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
			return
		}
		if err := tmpl.Execute(w, nil); err != nil {
			log.Println("Erreur d'exécution du template:", err) // Log d'erreur pour l'exécution
			http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
			return
		}
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

		// Recherche de l'utilisateur dans la base de données
		var storedUser User
		err := DB.QueryRow(`SELECT id, password FROM users WHERE email = ?`, user.Email).
			Scan(&storedUser.ID, &storedUser.Password)
		if err == sql.ErrNoRows {
			log.Println("Utilisateur non trouvé pour l'email:", user.Email) // Log si utilisateur non trouvé
			http.Error(w, "Email ou mot de passe incorrect", http.StatusUnauthorized)
			return
		} else if err != nil {
			log.Println("Erreur SQL:", err) // Log pour les autres erreurs SQL
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}

		// Vérification du mot de passe
		if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)); err != nil {
			log.Println("Mot de passe incorrect pour l'email:", user.Email) // Log pour mot de passe incorrect
			http.Error(w, "Email ou mot de passe incorrect", http.StatusUnauthorized)
			return
		}

		// Création de la session
		sessionID := CreateSession(storedUser.ID)
		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    sessionID,
			HttpOnly: true,
			Secure:   true, // Activer en HTTPS
			SameSite: http.SameSiteStrictMode,
		})

		// Réponse JSON de connexion réussie
		json.NewEncoder(w).Encode(map[string]string{"message": "Connexion réussie"})
	}
}
