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

func HomeHandler(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "session-name")
    email, ok := session.Values["email"].(string)

    data := struct {
        Email string
        Posts []Post
    }{}

    if ok && email != "" {
        data.Email = email
    }
    if DB == nil {
        log.Println("La connexion à la base de données est nulle")
    }
    rows, err := DB.Query("SELECT id, title, content, image, category FROM posts")
    if err != nil {
        log.Println("Erreur SQL lors de la récupération des posts :", err)
        http.Error(w, "Erreur lors de la récupération des posts", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var posts []Post
    for rows.Next() {
        var post Post
        if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Image, &post.Category); err != nil {
            log.Println("Erreur lors de la lecture des posts :", err)
            http.Error(w, "Erreur lors de la lecture des posts", http.StatusInternalServerError)
            return
        }
        posts = append(posts, post)
    }

    data.Posts = posts

    tmpl, err := template.ParseFiles("web/html/home.html")
    if err != nil {
        log.Println("Erreur lors du chargement du template :", err)
        http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
        return
    }

    if err := tmpl.Execute(w, data); err != nil {
        log.Println("Erreur lors de l'exécution du template :", err)
        http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
    }
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

        query := `INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)`
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

    session, err := store.Get(r, "session-name")
    if err != nil {
        http.Error(w, "Erreur de session", http.StatusInternalServerError)
        return
    }
    session.Values["email"] = loginReq.Email
    session.Save(r, w)

    log.Println("Connexion réussie pour l'utilisateur :", loginReq.Email)

    http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "session-name")
    session.Options.MaxAge = -1
    session.Save(r, w)

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"redirect": "/home"})
}

// Page de création de postfunc PostHandler(w http.ResponseWriter, r *http.Request) {
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
    
            // Insertion du post dans la base de données
            query := `INSERT INTO posts (id, title, content, image, category) VALUES (?, ?, ?, ?, ?)`
            _, err := DB.Exec(query, post.ID, post.Title, post.Content, post.Image, post.Category)
            if err != nil {
                log.Println("Erreur lors de l'insertion du post :", err)
                http.Error(w, "Erreur lors de la création du post", http.StatusInternalServerError)
                return
            }
    
            log.Println("Post inséré avec succès :", post.Title)
    
            // Redirection vers /home après la création du post
            w.WriteHeader(http.StatusCreated)
            json.NewEncoder(w).Encode(map[string]string{"redirect": "/home"})        }
    }
    