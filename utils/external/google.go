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
	"golang.org/x/oauth2/google"
)

// GoogleOauthConfig stocke la configuration OAuth pour Google
var GoogleOauthConfig *oauth2.Config

// InitGoogleOAuth initialise la configuration OAuth de Google
// Cette fonction charge les identifiants depuis le fichier external.env
func InitGoogleOAuth() {
	// Chargement des variables d'environnement
	err := loadEnv("./external.env")
	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <google.go> Impossible d'ouvrir le fichier env. Vérifiez que le fichier existe", err)
		logs.AddLogsToDatabase(logMsg)
	}

	// Configuration du client OAuth avec les identifiants Google
	GoogleOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:5080/auth/google/callback", // URL de redirection après autorisation
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",   // Permission pour accéder à l'email
			"https://www.googleapis.com/auth/userinfo.profile", // Permission pour accéder au profil (nom, photo)
		},
		Endpoint: google.Endpoint, // Utilise les endpoints OAuth officiels de Google
	}
}

// HandleGoogleLogin redirige l'utilisateur vers la page de consentement Google
// C'est la première étape du processus OAuth : demander l'autorisation à l'utilisateur
func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := GoogleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// HandleGoogleCallback gère la redirection de retour depuis Google après autorisation
// C'est ici que l'on traite la réponse de Google et qu'on crée/connecte l'utilisateur
func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	// ÉTAPE 1 : Récupération du code d'autorisation depuis l'URL
	code := r.URL.Query().Get("code")
	if code == "" {
		logMsg := "ERREUR : <google.go> Erreur dans la tentative de connexion, Google n'a pas renvoyé de code d'autorisation."
		logs.AddLogsToDatabase(logMsg)
		utils.StatusBadRequest(w)
		return
	}

	// ÉTAPE 2 : Échange du code d'autorisation contre un token d'accès
	token, err := GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <google.go> Erreur dans l'utilisation du code d'autorisation : ", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}

	// ÉTAPE 3 : Récupération des informations utilisateur via l'API Google
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <google.go> Impossible de récupérer les données de l'utilisateur : ", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}
	defer resp.Body.Close()

	// Décodage de la réponse JSON contenant les informations utilisateur
	var userInfo map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&userInfo)

	// ÉTAPE 4 : Extraction et validation des données essentielles
	// Récupération de l'ID Google (identifiant unique de l'utilisateur chez Google)
	googleID, ok := userInfo["id"].(string)
	if !ok {
		logMsg := "ERREUR : <google.go> ID utilisateur Google manquant"
		logs.AddLogsToDatabase(logMsg)
		return
	}

	// Récupération de l'email (obligatoire pour notre système)
	email, ok := userInfo["email"].(string)
	if !ok {
		logMsg := "ERREUR : <google.go> Email utilisateur Google manquant"
		logs.AddLogsToDatabase(logMsg)
		return
	}

	// Récupération du nom (optionnel, avec valeur par défaut)
	googleName, ok := userInfo["name"].(string)
	if !ok {
		googleName = "GoogleUser" // Nom par défaut si non fourni
	}

	// ÉTAPE 5 : Recherche ou création de l'utilisateur dans la base de données locale
	userID, err := GoogleUser(googleID, email, googleName)
	if err != nil {
		logMsg := fmt.Sprint("Échec de la recherche/création de l'utilisateur : ", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}

	// ÉTAPE 6 : Création de la session utilisateur (cookie)
	err = authhandlers.InitSession(w, userID, "user", googleName)
	if err != nil {
		utils.InternalServError(w)
		return
	}

	// ÉTAPE 7 : Redirection vers la page d'accueil
	http.Redirect(w, r, "/", http.StatusFound)
}

// GoogleUser gère la logique de recherche ou de création d'un utilisateur dans la base de données
func GoogleUser(googleID, email, username string) (int, error) {
	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		return 0, err
	}
	defer db.Close()

	var userID int

	// CAS 1 : Recherche d'un utilisateur ayant déjà ce google_id
	sqlQuery := `SELECT id FROM user WHERE google_id = ?`
	row := db.QueryRow(sqlQuery, googleID)
	err = row.Scan(&userID)

	if err == nil {
		// L'utilisateur a été trouvé avec ce google_id, on renvoie son ID pour le connecter
		return userID, nil
	} else if err != sql.ErrNoRows {
		// Erreur inattendue dans la base de données
		return 0, err
	}

	// CAS 2 et 3 : L'utilisateur n'a pas lié son compte Google
	// On vérifie s'il n'a pas créé un compte classique avec cette adresse mail
	if err == sql.ErrNoRows {
		// Recherche d'un utilisateur avec cette adresse email
		sqlQuery = `SELECT id FROM user WHERE email = ?`
		row = db.QueryRow(sqlQuery, email)
		err = row.Scan(&userID)

		switch err {
		// CAS 2 : L'utilisateur existe avec cet email → on associe son google_id
		// Cela permet à l'utilisateur de se connecter via Google à l'avenir
		case nil:
			sqlUpdate := `UPDATE user SET google_id = ? WHERE id = ?`
			_, err = db.Exec(sqlUpdate, googleID, userID)
			if err != nil {
				return 0, err
			}
		// CAS 3 : Aucun utilisateur n'existe avec cette adresse mail ou ce google_id
		// On crée un nouveau compte dans la base de données
		case sql.ErrNoRows:
			userID, err = CreateNewGoogleUser(googleID, email, username, db)
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

// CreateNewGoogleUser crée un nouvel utilisateur dans la base de données avec ses informations Google
func CreateNewGoogleUser(googleID, email, googleName string, db *sql.DB) (int, error) {
	// ÉTAPE 1 : Détermination du rôle de l'utilisateur
	var count int
	role := 3 // Rôle par défaut (simple membre)

	// Compte le nombre total d'utilisateurs dans la base
	err := db.QueryRow("SELECT COUNT(*) FROM user").Scan(&count)
	if err != nil {
		return 0, err
	}

	// Le premier utilisateur à s'inscrire devient automatiquement administrateur
	if count == 0 {
		role = 1
	}

	// ÉTAPE 2 : Génération d'un nom d'utilisateur unique
	// Si le nom est déjà pris, on ajoute un suffixe numérique (_1, _2, _3, etc.)
	addon := 0
	for {
		var id int
		testedName := googleName
		if addon != 0 {
			// Construction du nom avec suffixe : nom_1, nom_2, etc.
			testedName = fmt.Sprintf("%s_%d", googleName, addon)
		}

		// Vérification si ce nom d'utilisateur existe déjà
		sqlQuery := `SELECT id FROM user WHERE username = ?`
		row := db.QueryRow(sqlQuery, testedName)
		err = row.Scan(&id)

		if err != sql.ErrNoRows {
			if err == nil {
				// Le nom existe déjà, on incrémente le suffixe et on réessaie
				addon += 1
				continue
			} else {
				// Erreur de base de données
				return 0, err
			}
		} else {
			// Le nom est disponible, on l'utilise
			googleName = testedName
			break
		}
	}

	// ÉTAPE 3 : Insertion du nouvel utilisateur dans la base de données
	// Note : lae table 'user' a une nouvelle colonne 'google_id'
	sqlUpdate := `INSERT INTO user(username, email, google_id, role_id) VALUES(?, ?, ?, ?)`
	result, err := db.Exec(sqlUpdate, googleName, email, googleID, role)
	if err != nil {
		return 0, err
	}

	// Récupération de l'ID du nouvel utilisateur créé
	// Cet ID sera utilisé pour la session
	userID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(userID), nil
}
