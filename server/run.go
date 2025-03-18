package engine

import (
    "log"
    "net/http"
)

func Run(forum *User) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.Redirect(w, r, "/Home", http.StatusFound)
    })
	http.HandleFunc("/Home", HomeHandler)
	http.HandleFunc("/login", LoginHandler)
    http.HandleFunc("/register", RegisterHandler)


    log.Println("Serveur lancé sur http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}