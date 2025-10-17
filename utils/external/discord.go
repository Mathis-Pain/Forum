package external

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Mathis-Pain/Forum/handlers/authhandlers" // Assuming the path is correct
	"github.com/Mathis-Pain/Forum/utils"                 // Assuming the path is correct
	"golang.org/x/oauth2"
)

// Discord OAuth Endpoint
var DiscordEndpoint = oauth2.Endpoint{
	AuthURL:  "https://discord.com/api/oauth2/authorize",
	TokenURL: "https://discord.com/api/oauth2/token",
}

var DiscordOauthConfig *oauth2.Config

func InitDiscordOAuth() {
	err := loadEnv("./external.env")
	if err != nil {
		log.Print("Erreur à l'ouverture du fichier env pour Discord:", err)
	}

	DiscordOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("DISCORD_CLIENT_ID"),
		ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:5080/auth/discord/callback",
		Scopes: []string{
			"identify",
			"email",
		},
		Endpoint: DiscordEndpoint,
	}
}

func HandleDiscordLogin(w http.ResponseWriter, r *http.Request) {
	if DiscordOauthConfig == nil {
		http.Error(w, "Discord configuration not initialized.", http.StatusInternalServerError)
		return
	}
	url := DiscordOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleDiscordCallback(w http.ResponseWriter, r *http.Request) {
	// 1. Exchange the code for a token
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code manquant dans l'URL", http.StatusBadRequest)
		return
	}

	token, err := DiscordOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Échec lors de l'échange du code : "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 2. Fetch user information
	// We use the @me endpoint for both identify and email scopes
	resp, err := http.Get("https://discord.com/api/v10/users/@me?access_token=" + token.AccessToken)
	if err != nil {
		http.Error(w, "Impossible de récupérer les infos utilisateur", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&userInfo)

	// 3. Extract user data
	discordID, ok := userInfo["id"].(string)
	if !ok {
		http.Error(w, "ID utilisateur Discord manquant", http.StatusInternalServerError)
		return
	}

	email, ok := userInfo["email"].(string)
	if !ok || email == "" {
		http.Error(w, "Email utilisateur Discord manquant/non autorisé", http.StatusInternalServerError)
		return
	}

	// Discord combines username and discriminator (#1234) for a unique name.
	// Use the global name for a cleaner display, or build the old format.
	username, ok := userInfo["username"].(string)
	if !ok {
		username = "DiscordUser"
	}

	// 4. Find/Create user in local DB
	userID, err := DiscordUser(discordID, email, username)
	if err != nil {
		http.Error(w, "Échec de la recherche/création de l'utilisateur local: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 5. Create session and redirect
	err = authhandlers.InitSession(w, userID, "user", username)
	if err != nil {
		utils.InternalServError(w)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

// DiscordUser and CreateNewDiscordUser are the database logic functions.
// They need to be in a shared database/user package or defined here,
// similar to the DiscordUser/CreateNewDiscordUser functions from the previous step.
//
// ⚠️ IMPORTANT: You must modify your 'user' database table to add a 'discord_id' column.
// For brevity, the full database functions (DiscordUser, CreateNewDiscordUser)
// are omitted here but follow the exact logic of your existing Google/discord ones,
// just swapping 'google_id'/'discord_id' for 'discord_id'.
func DiscordUser(discordID, email, username string) (int, error) {
	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		return 0, err
	}
	defer db.Close()

	var userID int

	// Cherche l'utilisateur ayant ce discord_id dans la base de données
	sqlQuery := `SELECT id FROM user WHERE discord_id = ?`
	row := db.QueryRow(sqlQuery, discordID)
	err = row.Scan(&userID)

	if err == nil {
		// L'utilisateur a été trouvé, renvoie son id
		return userID, nil
	} else if err != sql.ErrNoRows {
		// Erreur dans la base de données
		return 0, err
	}

	// L'utilisateur n'a pas lié son compte discord, on vérifie s'il n'a pas utilisé cette adresse mail
	if err == sql.ErrNoRows {
		sqlQuery = `SELECT id FROM user WHERE email = ?`
		row = db.QueryRow(sqlQuery, email)
		err = row.Scan(&userID)

		switch err {
		// L'utilisateur a été trouvé, on associe son discord_id à son adresse mail
		case nil:
			sqlUpdate := `UPDATE user SET discord_id = ? WHERE id = ?`
			_, err = db.Exec(sqlUpdate, discordID, userID)
			if err != nil {
				return 0, err
			}
		// Aucun utilisateur n'existe, on l'ajoute
		case sql.ErrNoRows:
			// You'll need to define this function or move it to a shared package
			userID, err = CreateNewDiscordUser(discordID, email, username, db)
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

func CreateNewDiscordUser(discordID, email, discordName string, db *sql.DB) (int, error) {
	// This is essentially the same logic as CreateNewGoogleUser but for discord.
	// It should handle role assignment, unique username creation, and insertion.

	// 1. Determine role
	var count int
	role := 3 // Default role
	err := db.QueryRow("SELECT COUNT(*) FROM user").Scan(&count)
	if err != nil {
		return 0, err
	}
	if count == 0 {
		role = 1 // Admin for first user
	}

	// 2. Ensure unique username
	addon := 0
	uniqueUsername := discordName
	for {
		var id int
		testedName := discordName
		if addon != 0 {
			testedName = fmt.Sprintf("%s_%d", discordName, addon)
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
			uniqueUsername = testedName
			break
		}
	}

	// 3. Insert new user
	// Note: You must ensure your 'user' table has a 'discord_id' column.
	sqlUpdate := `INSERT INTO user(username, email, discord_id, role_id) VALUES(?, ?, ?, ?)`
	result, err := db.Exec(sqlUpdate, uniqueUsername, email, discordID, role)
	if err != nil {
		return 0, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(userID), nil
}
