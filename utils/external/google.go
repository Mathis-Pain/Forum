package external

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Mathis-Pain/Forum/handlers/authhandlers"
	"github.com/Mathis-Pain/Forum/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var GoogleOauthConfig *oauth2.Config

func loadEnv(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		// It's often non-fatal to not find an env file, depending on your deployment.
		// You might change this to log and return nil if you expect envs to be set externally.
		return fmt.Errorf("error opening .env file %s: %w", filename, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split only on the first '='
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// Set the environment variable
			os.Setenv(key, value)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading .env file: %w", err)
	}

	return nil
}

func InitGoogleOAuth() {
	err := loadEnv("./google.env")
	if err != nil {
		log.Print("Erreur à l'ouverture du fichier env :", err)
	}

	GoogleOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:5080/auth/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := GoogleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Récupération des données utilisateur transmises par google
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code manquant dans l'URL", http.StatusBadRequest)
		return
	}

	token, err := GoogleOauthConfig.Exchange(context.Background(), code)
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

	// Enregistrement des données pour la recherche ou la création du compte
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
		googleName = "GoogleUser"
	}

	userID, err := GoogleUser(googleID, email, googleName)
	if err != nil {
		http.Error(w, "Échec de la recherche/création de l'utilisateur local: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Création du cookie
	err = authhandlers.InitSession(w, userID, "user", googleName)
	if err != nil {
		utils.InternalServError(w)
		return
	}

	// Redirection
	http.Redirect(w, r, "/", http.StatusFound)
}

func GoogleUser(googleID, email, username string) (int, error) {
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

	sqlUpdate := `INSERT INTO user(username, email, google_id, role_id) VALUES(?, ?, ?, ?, ?)`
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
