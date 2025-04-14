package engine

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
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
        tmpl, err := template.ParseFiles("web/html/register.html")
        if err != nil {
            http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
            return
        }
        tmpl.Execute(w, nil)
        return
    }

    if r.Method == "POST" {
        body, _ := io.ReadAll(r.Body)
        log.Println("Corps de la requête reçu :", string(body)) // Log des données reçues

        var user User
        if err := json.Unmarshal(body, &user); err != nil {
            log.Println("Erreur de décodage JSON :", err)
            http.Error(w, "Données invalides", http.StatusBadRequest)
            return
        }

        log.Println("Données utilisateur après décodage :", user)

        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
        if err != nil {
            log.Println("Erreur lors du hachage du mot de passe :", err)
            http.Error(w, "Erreur serveur", http.StatusInternalServerError)
            return
        }

        query := `INSERT INTO users (username, email, password) VALUES (?, ?, ?)`
        _, err = DB.Exec(query, user.Username, user.Email, hashedPassword)
        if err != nil {
            log.Println("Erreur lors de l'insertion utilisateur :", err)
            http.Error(w, "Erreur lors de l'inscription", http.StatusInternalServerError)
            return
        }

        log.Println("Utilisateur inséré avec succès :", user.Username)
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(map[string]string{"message": "Inscription réussie"})
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

    if r.Method != "POST" {
        http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
        return
    }

    var loginReq User
    err := json.NewDecoder(r.Body).Decode(&loginReq)
    if err != nil {
        log.Println("Erreur de décodage JSON :", err)
        http.Error(w, "Données invalides", http.StatusBadRequest)
        return
    }

    var storedPassword string
    var userID int
    err = DB.QueryRow(`SELECT id, password_hash FROM users WHERE email = ?`, loginReq.Email).Scan(&userID, &storedPassword)
    if err == sql.ErrNoRows {
        http.Error(w, "Utilisateur non trouvé", http.StatusUnauthorized)
        return
    } else if err != nil {
        log.Println("Erreur SQL :", err)
        http.Error(w, "Erreur interne", http.StatusInternalServerError)
        return
    }

    err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(loginReq.Password))
    if err != nil {
        http.Error(w, "Mot de passe incorrect", http.StatusUnauthorized)
        return
    }

    // Connexion réussie
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Connexion réussie"})
}

func AccueilHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        tmpl, err := template.ParseFiles("web/html/accueil.html")
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
	}
}

// Page de création de post
func PostHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        tmpl, err := template.ParseFiles("web/html/post.html")
        if err != nil {
            log.Println("Erreur lors du chargement du template :", err)
            http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
            return
        }
        tmpl.Execute(w, nil)
        return
    }

    if r.Method == "POST" {
        body, _ := io.ReadAll(r.Body)
        log.Println("Corps de la requête reçu :", string(body)) // Log des données reçues

        var post Post
        if err := json.Unmarshal(body, &post); err != nil {
            log.Println("Erreur de décodage JSON :", err)
            http.Error(w, "Données invalides", http.StatusBadRequest)
            return
        }

        log.Println("Données post après décodage :", post)

        query := `INSERT INTO posts (user_id, title, content, image, category) VALUES (?, ?, ?, ?, ?)`
        _, err := DB.Exec(query, post.UserID, post.Title, post.Content, post.Image, post.Category)
        if err != nil {
            log.Println("Erreur lors de l'insertion du post :", err)
            http.Error(w, "Erreur lors de la création du post", http.StatusInternalServerError)
            return
        }

        log.Println("Post inséré avec succès :", post.Title)
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(map[string]string{"message": "Post créé avec succès"})
    }
}

// Menu pour voir les post
func MenuHandler(w http.ResponseWriter, r *http.Request) {
    
}