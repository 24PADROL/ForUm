package engine

import (
    "database/sql"
    "encoding/json"
    "html/template"
    "log"
    "net/http"
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
        err := json.NewDecoder(r.Body).Decode(&user)
        if err != nil {
            log.Println("Erreur de décodage JSON:", err)
            http.Error(w, "Données invalides", http.StatusBadRequest)
            return
        }

        log.Println("Données utilisateur reçues:", user)

        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
        if err != nil {
            log.Println("Erreur lors du hachage du mot de passe:", err)
            http.Error(w, "Erreur interne", http.StatusInternalServerError)
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

    if r.Method == "POST" {
        var user User
        err := json.NewDecoder(r.Body).Decode(&user)
        if err != nil {
            log.Println("Erreur de décodage JSON:", err)
            http.Error(w, "Données invalides", http.StatusBadRequest)
            return
        }

        var storedUser User
        err = DB.QueryRow(`SELECT id, password FROM users WHERE email = ?`, user.Email).Scan(&storedUser.ID, &storedUser.Password)
        if err == sql.ErrNoRows {
            http.Error(w, "Utilisateur non trouvé", http.StatusUnauthorized)
            return
        } else if err != nil {
            log.Println("Erreur SQL:", err)
            http.Error(w, "Erreur serveur", http.StatusInternalServerError)
            return
        }

        err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
        if err != nil {
            http.Error(w, "Mot de passe incorrect", http.StatusUnauthorized)
            return
        }

        sessionID := CreateSession(storedUser.ID)
        http.SetCookie(w, &http.Cookie{Name: "session", Value: sessionID, HttpOnly: true})

        // Rediriger vers la page d'accueil
        http.Redirect(w, r, "/accueil", http.StatusSeeOther)
    }
}

func AccueilHandler(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("web/html/accueil.html"))
    tmpl.Execute(w, nil)
}