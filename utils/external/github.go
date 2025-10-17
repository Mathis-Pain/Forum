package external

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Mathis-Pain/Forum/handlers/authhandlers"
	"github.com/Mathis-Pain/Forum/utils"
	"github.com/Mathis-Pain/Forum/utils/logs"
	"golang.org/x/oauth2"
)

// GitHubEndpoint définit les URLs nécessaires pour l'authentification OAuth avec GitHub
// GitHubOauthConfig stocke la configuration OAuth pour GitHub
var GitHubEndpoint = oauth2.Endpoint{
	AuthURL:  "https://github.com/login/oauth/authorize",    // URL pour demander l'autorisation
	TokenURL: "https://github.com/login/oauth/access_token", // URL pour échanger le code contre un token
}

var GitHubOauthConfig *oauth2.Config

// Appelée dans le main, InitGitHubOAuth initialise la configuration OAuth de GitHub
// Cette fonction charge les identifiants depuis le fichier external.env
func InitGitHubOAuth() {
	// Chargement des variables d'environnement depuis le fichier external.env
	err := loadEnv("./external.env")
	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <github.go> Impossible d'ouvrir le fichier env. Vérifiez que le fichier existe", err)
		logs.AddLogsToDatabase(logMsg)
	}

	// Récupère les identifiants GitHub dans le .env
	GitHubOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:5080/auth/github/callback", // URL de redirection configurée sur GitHub et dans les routes
		Scopes: []string{
			"user:email",
		},
		Endpoint: GitHubEndpoint,
	}
}

// HandleGitHubLogin redirige l'utilisateur vers la page de consentement GitHub
// C'est la première étape du processus OAuth : demander l'autorisation à l'utilisateur
func HandleGitHubLogin(w http.ResponseWriter, r *http.Request) {
	url := GitHubOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// HandleGitHubCallback gère la redirection après autorisation
// C'est ici que l'on traite la réponse de GitHub et qu'on crée/connecte l'utilisateur
func HandleGitHubCallback(w http.ResponseWriter, r *http.Request) {
	// Récupération du code d'autorisation depuis l'URL
	code := r.URL.Query().Get("code")
	if code == "" {
		logMsg := "ERREUR : <github.go> Erreur dans la tentative de connexion, GitHub n'a pas renvoyé de code d'autorisation."
		logs.AddLogsToDatabase(logMsg)
		utils.StatusBadRequest(w)
		return
	}

	// Échange du code d'autorisation contre un token d'accès
	token, err := GitHubOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <github.go> Erreur dans l'utilisateur du code d'autorisation : ", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}

	// ÉTAPE 1 : Récupération des informations de l'utilisateur (ID, login)
	req, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	req.Header.Set("Authorization", "token "+token.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <github.go> Impossible de récupérer les données de l'utilisateur : ", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}
	defer resp.Body.Close()

	// Décodage de la réponse JSON contenant les informations utilisateur
	var userInfo map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&userInfo)

	// Conversion de l'ID GitHub (nombre) en chaîne pour cohérence avec le stockage en base
	githubID := fmt.Sprintf("%.0f", userInfo["id"].(float64))
	// Récupération du nom d'utilisateur GitHub
	githubUsername, ok := userInfo["login"].(string)
	if !ok {
		githubUsername = "GitHubUser" // Valeur par défaut si le login n'est pas disponible
	}

	// ÉTAPE 2 : Récupération de l'email principal
	// Un appel séparé est nécessaire car l'email peut être privé ou null dans l'API de base
	email, err := getGitHubPrimaryEmail(token.AccessToken)
	if err != nil || email == "" {
		logMsg := fmt.Sprint("ERREUR : <github.go> Erreur dans la récupération de l'email : ", err, " création d'un mail placeholder pour la base de données.")
		logs.AddLogsToDatabase(logMsg)

		// Génération d'un email de secours pour la base de données
		email = fmt.Sprintf("%s@github-user.noemail", githubID)
	}

	// ÉTAPE 3 : Recherche ou création de l'utilisateur dans la base de données locale
	userID, err := GitHubUser(githubID, email, githubUsername)
	if err != nil {
		logMsg := fmt.Sprint("Échec de la recherche/création de l'utilisateur : ", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}

	// ÉTAPE 4 : Création de la session utilisateur (cookie)
	err = authhandlers.InitSession(w, userID, "user", githubUsername)
	if err != nil {
		utils.InternalServError(w)
		return
	}

	// ÉTAPE 5 : Redirection vers la page d'accueil
	http.Redirect(w, r, "/", http.StatusFound)
}

// getGitHubPrimaryEmail récupère l'email principal et vérifié de l'utilisateur GitHub
// Cette fonction est nécessaire car l'email n'est pas toujours disponible dans l'API de base
func getGitHubPrimaryEmail(accessToken string) (string, error) {
	// Requête vers l'endpoint des emails GitHub
	req, _ := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	req.Header.Set("Authorization", "token "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Décodage de la liste des emails
	var emails []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", err
	}

	// Recherche de l'email principal et vérifié
	for _, e := range emails {
		isPrimary, ok1 := e["primary"].(bool)
		isVerified, ok2 := e["verified"].(bool)
		email, ok3 := e["email"].(string)

		// Retourne le premier email qui est à la fois principal et vérifié
		if ok1 && ok2 && ok3 && isPrimary && isVerified {
			return email, nil
		}
	}
	return "", fmt.Errorf("aucun email principal et vérifié trouvé")
}

// GitHubUser gère la logique de recherche ou de création d'un utilisateur dans la base de données locale
func GitHubUser(githubID, email, username string) (int, error) {
	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		return 0, err
	}
	defer db.Close()

	var userID int

	// CAS 1 : Recherche d'un utilisateur ayant déjà ce github_id
	sqlQuery := `SELECT id FROM user WHERE github_id = ?`
	row := db.QueryRow(sqlQuery, githubID)
	err = row.Scan(&userID)

	if err == nil {
		// L'utilisateur a été trouvé avec ce github_id, on renvoie son ID
		return userID, nil
	} else if err != sql.ErrNoRows {
		// Erreur inattendue dans la base de données
		return 0, err
	}

	// CAS 2 et 3 : L'utilisateur n'a pas lié son compte GitHub
	if err == sql.ErrNoRows {
		// Recherche d'un utilisateur avec cette adresse email
		sqlQuery = `SELECT id FROM user WHERE email = ?`
		row = db.QueryRow(sqlQuery, email)
		err = row.Scan(&userID)

		switch err {
		// CAS 2 : L'utilisateur existe avec cet email → on associe son github_id
		case nil:
			sqlUpdate := `UPDATE user SET github_id = ? WHERE id = ?`
			_, err = db.Exec(sqlUpdate, githubID, userID)
			if err != nil {
				return 0, err
			}
		// CAS 3 : Aucun utilisateur n'existe → on crée un nouveau compte
		case sql.ErrNoRows:
			userID, err = CreateNewGitHubUser(githubID, email, username, db)
			if err != nil {
				return 0, err
			}
		default:
			// Erreur inattendue dans la base de données
			return 0, err
		}
	}

	return userID, nil
}

// CreateNewGitHubUser crée un nouvel utilisateur dans la base de données avec ses informations GitHub
// Cette fonction gère l'attribution du rôle, la création d'un nom d'utilisateur unique et l'insertion en base
func CreateNewGitHubUser(githubID, email, githubName string, db *sql.DB) (int, error) {
	// ÉTAPE 1 : Détermination du rôle de l'utilisateur
	var count int
	role := 3 // Rôle par défaut (simple membre)
	err := db.QueryRow("SELECT COUNT(*) FROM user").Scan(&count)
	if err != nil {
		return 0, err
	}
	// Le premier utilisateur devient administrateur
	if count == 0 {
		role = 1
	}

	// ÉTAPE 2 : Génération d'un nom d'utilisateur unique
	// Si le nom est déjà pris, on ajoute un suffixe numérique (_1, _2, etc.)
	addon := 0
	uniqueUsername := githubName
	for {
		var id int
		testedName := githubName
		if addon != 0 {
			testedName = fmt.Sprintf("%s_%d", githubName, addon)
		}
		// Vérifie si le nom d'utilisateur existe déjà
		sqlQuery := `SELECT id FROM user WHERE username = ?`
		row := db.QueryRow(sqlQuery, testedName)
		err = row.Scan(&id)
		if err != sql.ErrNoRows {
			if err == nil {
				// Le nom existe déjà, on incrémente le suffixe
				addon += 1
				continue
			} else {
				// Erreur de base de données
				return 0, err
			}
		} else {
			// Le nom est disponible
			uniqueUsername = testedName
			break
		}
	}

	// ÉTAPE 3 : Insertion du nouvel utilisateur dans la base de données
	// Note : la table 'user' contient une nouvelle colonne 'github_id'
	sqlUpdate := `INSERT INTO user(username, email, github_id, role_id) VALUES(?, ?, ?, ?)`
	result, err := db.Exec(sqlUpdate, uniqueUsername, email, githubID, role)
	if err != nil {
		return 0, err
	}

	// Récupération de l'ID du nouvel utilisateur créé
	userID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(userID), nil
}
