package engine

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "os"

	"github.com/gorilla/sessions"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
)

var googleOauthConfig *oauth2.Config
var store = sessions.NewCookieStore([]byte("client_secret_717864681118-74dt6ue2a8ilo1g7upk6m4gdqp2lmv5e.apps.googleusercontent.com.json"))
func init() {
    file, err := os.Open("client_secret_717864681118-74dt6ue2a8ilo1g7upk6m4gdqp2lmv5e.apps.googleusercontent.com.json")
    if err != nil {
        log.Fatalf("Erreur lors de l'ouverture du fichier JSON : %v", err)
    }
    defer file.Close()

    // Décoder le fichier JSON
    var config struct {
        Web struct {
            ClientID     string   `json:"client_id"`
            ClientSecret string   `json:"client_secret"`
            RedirectURIs []string `json:"redirect_uris"`
        } `json:"web"`
    }
    if err := json.NewDecoder(file).Decode(&config); err != nil {
        log.Fatalf("Erreur lors du décodage du fichier JSON : %v", err)
    }

    googleOauthConfig = &oauth2.Config{
        ClientID:     config.Web.ClientID,
        ClientSecret: config.Web.ClientSecret,
        RedirectURL:  config.Web.RedirectURIs[0],
        Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
        Endpoint:     google.Endpoint,
    }
}

func GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
    url := googleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
    http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
    code := r.URL.Query().Get("code")
    if code == "" {
        http.Error(w, "Code non trouvé", http.StatusBadRequest)
        return
    }

    token, err := googleOauthConfig.Exchange(context.Background(), code)
    if err != nil {
        log.Println("Erreur lors de l'échange de token :", err)
        http.Error(w, "Erreur lors de l'échange de token", http.StatusInternalServerError)
        return
    }

    client := googleOauthConfig.Client(context.Background(), token)
    resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
    if err != nil {
        log.Println("Erreur lors de la récupération des informations utilisateur :", err)
        http.Error(w, "Erreur lors de la récupération des informations utilisateur", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

	var userInfo struct {
        Email string `json:"email"`
        Name  string `json:"name"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
        log.Println("Erreur lors du décodage des informations utilisateur :", err)
        http.Error(w, "Erreur lors du décodage des informations utilisateur", http.StatusInternalServerError)
        return
    }

	session, _ := store.Get(r, "session-name")
    session.Values["email"] = userInfo.Email
    session.Save(r, w)

    http.Redirect(w, r, "/home", http.StatusSeeOther)
}