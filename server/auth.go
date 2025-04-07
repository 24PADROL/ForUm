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
	tmpl := template.Must(template.ParseFiles("template/Home.html"))
	tmpl.Execute(w, nil)
}

// Inscription
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("template/register.html")
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

func MessagesHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer les messages depuis la base de données
	messages, err := GetMessages()
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des messages", http.StatusInternalServerError)
		return
	}

	// Rendre la page HTML avec les messages
	tmpl, err := template.ParseFiles("template/messages.html")
	if err != nil {
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, messages)
}



// Page de discussion
func ForumHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        // Récupérer les messages depuis la base de données
        rows, err := DB.Query("SELECT id, username, message FROM posts ORDER BY created_at DESC")
        if err != nil {
            http.Error(w, "Erreur lors de la récupération des messages", http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        var posts []Post
        for rows.Next() {
            var post Post
            if err := rows.Scan(&post.ID, &post.Username, &post.Message); err != nil {
                http.Error(w, "Erreur de lecture des messages", http.StatusInternalServerError)
                return
            }
            posts = append(posts, post)
        }

        // Rendu de la page avec les messages
        tmpl := template.Must(template.ParseFiles("template/forum.tmpl"))
        tmpl.Execute(w, posts)
    } else {
        http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
    }
}


func GetMessages() ([]Message, error) {
	// Récupérer les messages depuis la base de données
	rows, err := DB.Query("SELECT id, user_id, content, created_at FROM messages ORDER BY created_at DESC")
	if err != nil {
		log.Println("Erreur lors de la récupération des messages:", err)
		return nil, err
	}
	defer rows.Close()

	var messages []Message

	for rows.Next() {
		var message Message
		if err := rows.Scan(&message.ID, &message.UserID, &message.Content, &message.CreatedAt); err != nil {
			log.Println("Erreur lors du scan des résultats:", err)
			return nil, err
		}
		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		log.Println("Erreur lors de l'itération sur les lignes:", err)
		return nil, err
	}

	return messages, nil
}



// Poster un message
func PostMessageHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        var message Message
        // Lire les données du message envoyé en JSON
        if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
            http.Error(w, "Données invalides", http.StatusBadRequest)
            return
        }

        // Assurer qu'il y a un utilisateur connecté et que l'ID de l'utilisateur est valide
        userID := 1 // Remplace cela par la logique pour obtenir l'ID utilisateur actuel via la session ou autre.

        // Insérer le message dans la base de données
        err := InsertMessage(userID, message.Content)
        if err != nil {
            http.Error(w, "Erreur lors de la publication du message", http.StatusInternalServerError)
            return
        }

        // Répondre avec un message de succès
        json.NewEncoder(w).Encode(map[string]string{"message": "Message publié avec succès"})
    } else {
        http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
    }
}


func InsertMessage(userID int, content string) error {
	// Insérer un message dans la base de données
	_, err := DB.Exec("INSERT INTO messages (user_id, content) VALUES (?, ?)", userID, content)
	if err != nil {
		log.Println("Erreur lors de l'insertion du message:", err)
		return err
	}
	log.Println("Message ajouté avec succès !")
	return nil
}
