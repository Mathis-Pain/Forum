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

// GitHub OAuth Endpoint
var GitHubEndpoint = oauth2.Endpoint{
	AuthURL:  "https://github.com/login/oauth/authorize",
	TokenURL: "https://github.com/login/oauth/access_token",
}

var GitHubOauthConfig *oauth2.Config

// InitGitHubOAuth initializes the GitHub OAuth configuration
func InitGitHubOAuth() {
	// Use the same loadEnv but point to github.env
	err := loadEnv("./external.env")
	if err != nil {
		log.Print("Erreur à l'ouverture du fichier env pour GitHub:", err)
	}

	GitHubOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		// Must match the "Authorization callback URL" set on GitHub
		RedirectURL: "http://localhost:5080/auth/github/callback",
		Scopes: []string{
			"user:email", // Scope to get user's email
		},
		Endpoint: GitHubEndpoint,
	}
}

// HandleGitHubLogin redirects the user to GitHub's consent page
func HandleGitHubLogin(w http.ResponseWriter, r *http.Request) {
	url := GitHubOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// HandleGitHubCallback handles the redirect back from GitHub
func HandleGitHubCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code manquant dans l'URL", http.StatusBadRequest)
		return
	}

	token, err := GitHubOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Échec lors de l'échange du code : "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 1. Get basic user info (ID, login)
	req, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	req.Header.Set("Authorization", "token "+token.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Impossible de récupérer les infos utilisateur (base)", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&userInfo)

	// GitHub ID is typically a number, we convert it to string for consistency with google_id storage
	githubID := fmt.Sprintf("%.0f", userInfo["id"].(float64))
	githubUsername, ok := userInfo["login"].(string)
	if !ok {
		githubUsername = "GitHubUser"
	}

	// 2. Get email info (requires separate call as the basic 'user' endpoint might not expose the primary email)
	// This is required because the email can be null or private in the first API call.
	email, err := getGitHubPrimaryEmail(token.AccessToken)
	if err != nil || email == "" {
		// Fallback or error handling if email can't be retrieved
		log.Println("Could not retrieve primary email from GitHub:", err)
		// If the email is essential, you might stop here or ask the user to set a public email on GitHub.
		// For this example, we'll use a placeholder/generated email if it's strictly necessary for your DB.
		email = fmt.Sprintf("%s@github-user.noemail", githubID)
	}

	// Logic for finding/creating the user in your database
	userID, err := GitHubUser(githubID, email, githubUsername)
	if err != nil {
		http.Error(w, "Échec de la recherche/création de l'utilisateur local: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Création du cookie
	err = authhandlers.InitSession(w, userID, "user", githubUsername)
	if err != nil {
		utils.InternalServError(w)
		return
	}

	// Redirection
	http.Redirect(w, r, "/", http.StatusFound)
}

// getGitHubPrimaryEmail fetches the user's primary, verified email.
func getGitHubPrimaryEmail(accessToken string) (string, error) {
	req, _ := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	req.Header.Set("Authorization", "token "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var emails []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", err
	}

	for _, e := range emails {
		isPrimary, ok1 := e["primary"].(bool)
		isVerified, ok2 := e["verified"].(bool)
		email, ok3 := e["email"].(string)

		if ok1 && ok2 && ok3 && isPrimary && isVerified {
			return email, nil
		}
	}
	return "", fmt.Errorf("no primary and verified email found")
}

// GitHubUser handles the logic for finding or creating a user in the local database.
// This function needs to be adapted to handle a column like 'github_id' in your 'user' table.
// *You'll need to update your database schema to add a `github_id` column.*
func GitHubUser(githubID, email, username string) (int, error) {
	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		return 0, err
	}
	defer db.Close()

	var userID int

	// Cherche l'utilisateur ayant ce github_id dans la base de données
	sqlQuery := `SELECT id FROM user WHERE github_id = ?`
	row := db.QueryRow(sqlQuery, githubID)
	err = row.Scan(&userID)

	if err == nil {
		// L'utilisateur a été trouvé, renvoie son id
		return userID, nil
	} else if err != sql.ErrNoRows {
		// Erreur dans la base de données
		return 0, err
	}

	// L'utilisateur n'a pas lié son compte github, on vérifie s'il n'a pas utilisé cette adresse mail
	if err == sql.ErrNoRows {
		sqlQuery = `SELECT id FROM user WHERE email = ?`
		row = db.QueryRow(sqlQuery, email)
		err = row.Scan(&userID)

		switch err {
		// L'utilisateur a été trouvé, on associe son github_id à son adresse mail
		case nil:
			sqlUpdate := `UPDATE user SET github_id = ? WHERE id = ?`
			_, err = db.Exec(sqlUpdate, githubID, userID)
			if err != nil {
				return 0, err
			}
		// Aucun utilisateur n'existe, on l'ajoute
		case sql.ErrNoRows:
			// You'll need to define this function or move it to a shared package
			userID, err = CreateNewGitHubUser(githubID, email, username, db)
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

// CreateNewGitHubUser is a function placeholder, it should be the same as CreateNewGoogleUser
// but setting 'github_id' instead of 'google_id'. You might want to merge CreateNewGoogleUser
// and CreateNewGitHubUser into a single, generic user creation function.
func CreateNewGitHubUser(githubID, email, githubName string, db *sql.DB) (int, error) {
	// This is essentially the same logic as CreateNewGoogleUser but for GitHub.
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
	uniqueUsername := githubName
	for {
		var id int
		testedName := githubName
		if addon != 0 {
			testedName = fmt.Sprintf("%s_%d", githubName, addon)
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
	// Note: You must ensure your 'user' table has a 'github_id' column.
	sqlUpdate := `INSERT INTO user(username, email, github_id, role_id) VALUES(?, ?, ?, ?)`
	result, err := db.Exec(sqlUpdate, uniqueUsername, email, githubID, role)
	if err != nil {
		return 0, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(userID), nil
}
