package engine

import (
    "log"
    "net/http"

)

func Run(forum *User) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.Redirect(w, r, "/home", http.StatusFound)
    })
	http.HandleFunc("/home", HomeHandler)
	http.HandleFunc("/login", LoginHandler)
    http.HandleFunc("/register", RegisterHandler)
    http.HandleFunc("/accueil", AccueilHandler)
    http.HandleFunc("/post", PostHandler)
    http.HandleFunc("/auth/google", GoogleLoginHandler)
    http.HandleFunc("/auth/google/callback", GoogleCallbackHandler)
    http.HandleFunc("/logout", LogoutHandler)
    http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))

    log.Println("Serveur lancé sur http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}