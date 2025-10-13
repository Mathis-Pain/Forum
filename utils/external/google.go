package external

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Mathis-Pain/Forum/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig = &oauth2.Config{
		ClientID:     "ForumLocal",
		ClientSecret: "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",
		RedirectURL:  "http://localhost:5080/google/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
)

func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// external package (or wherever your OAuth handlers are)

func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	// ... (Existing code for exchanging code for token and fetching user info) ...

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code manquant dans l'URL", http.StatusBadRequest)
		return
	}

	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Échec lors de l'échange du code : "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		http.Error(w, "Impossible de récupérer les infos utilisateur", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&userInfo)

	// --- INTEGRATION STARTS HERE ---

	googleID, ok := userInfo["id"].(string)
	if !ok {
		http.Error(w, "ID utilisateur Google manquant", http.StatusInternalServerError)
		return
	}
	email, ok := userInfo["email"].(string)
	if !ok {
		http.Error(w, "Email utilisateur Google manquant", http.StatusInternalServerError)
		return
	}

	googleName, ok := userInfo["name"].(string)
	if !ok {
		// Fallback or use a generic name if 'name' is missing
		googleName = "GoogleUser"
	}

	// 1. Find or create the user in your local DB
	userID, err := FindOrCreateUserByGoogleID(googleID, email, googleName)
	if err != nil {
		http.Error(w, "Échec de la recherche/création de l'utilisateur local: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 2. Create the session
	session, err := sessions.CreateSession(userID)
	if err != nil {
		http.Error(w, "Échec de la création de la session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Set the session cookie
	// You should determine if you are running in a secure (HTTPS) environment
	isSecure := r.URL.Scheme == "https"
	sessions.SetCookie(w, "session_id", session.ID, isSecure)

	// 4. Redirect the user to a protected area or home page
	http.Redirect(w, r, "/", http.StatusFound)
}

func FindOrCreateUserByGoogleID(googleID, email, username string) (int, error) {
	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		return 0, err
	}
	defer db.Close()

	var userID int

	// Cherche l'utilisateur ayant ce google_id dans la base de données
	sqlQuery := `SELECT id FROM user WHERE google_id = ?`
	row := db.QueryRow(sqlQuery, googleID)
	err = row.Scan(&userID)

	if err == nil {
		// L'utilisateur a été trouvé, renvoie son id pour le connecter
		return userID, nil
	} else if err != sql.ErrNoRows {
		// Erreur dans la base de données
		return 0, err
	}

	// L'utilisateur n'a pas lié son compte google, on vérifie quand même s'il n'a pas utilisé cette adresse mail pour créer un compte classique
	if err == sql.ErrNoRows {
		sqlQuery = `SELECT id FROM user WHERE email = ?`
		row = db.QueryRow(sqlQuery, email)
		err = row.Scan(&userID)

		switch err {
		// L'utilisateur a été trouvé, on associe son google_id à son adresse mail pour qu'il puisse se connecter via google
		case nil:
			sqlUpdate := `UPDATE user SET google_id = ? WHERE id = ?`
			_, err = db.Exec(sqlUpdate, googleID, userID)
			if err != nil {
				return 0, err
			}
		// Aucun utilisateur n'existe avec cette adresse mail ou ce google_id, on l'ajoute à la base de données
		case sql.ErrNoRows:
			userID, err = CreateNewGoogleUser(googleID, email, username, db)
			if err != nil {
				return 0, err
			}
		default:
			// Erreur dans la base de données
			return 0, err
		}

	}

	return userID, nil
}

func CreateNewGoogleUser(googleID, email, googleName string, db *sql.DB) (int, error) {
	// ---- Vérifie si c'est le premier utilisateur ---
	var count int
	role := 3
	err := db.QueryRow("SELECT COUNT(*) FROM user").Scan(&count)
	if err != nil {
		return 0, err
	}
	if count == 0 {
		role = 1
	}

	// Vérifie si le nom d'utilisateur n'est pas déjà utilisé
	addon := 0
	for {
		var id int
		testedName := googleName
		if addon != 0 {
			testedName = fmt.Sprintf("%s_%d", googleName, addon)
		}
		sqlQuery := `SELECT id FROM user WHERE username = ?`
		row := db.QueryRow(sqlQuery, testedName)
		err = row.Scan(&id)
		if err != sql.ErrNoRows {
			if err == nil {
				addon += 1
				continue
			} else {
				return 0, err
			}
		} else {
			googleName = testedName
			break
		}
	}

	sqlUpdate := `INSERT INTO user(username, email, google_id, role_id) VALUES(?, ?, ?, ?)`
	result, err := db.Exec(sqlUpdate, googleName, email, googleID, role)
	if err != nil {
		return 0, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(userID), nil
}
