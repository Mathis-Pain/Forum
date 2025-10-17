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

// DiscordEndpoint définit les URLs nécessaires pour l'authentification OAuth avec Discord
var DiscordEndpoint = oauth2.Endpoint{
	AuthURL:  "https://discord.com/api/oauth2/authorize",
	TokenURL: "https://discord.com/api/oauth2/token",
}

// DiscordOauthConfig stocke la configuration OAuth pour Discord
var DiscordOauthConfig *oauth2.Config

// InitDiscordOAuth initialise la configuration OAuth de Discord
// Cette fonction charge les identifiants depuis le fichier external.env
func InitDiscordOAuth() {
	// Chargement des variables d'environnement depuis le fichier external.env
	err := loadEnv("./external.env")
	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <discord.go> Impossible d'ouvrir le fichier env. Vérifiez que le fichier existe", err)
		logs.AddLogsToDatabase(logMsg)
	}

	// Configuration du client OAuth avec les identifiants Discord
	DiscordOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("DISCORD_CLIENT_ID"),
		ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:5080/auth/discord/callback", // URL de redirection après autorisation
		Scopes: []string{
			"identify",
			"email",
		},
		Endpoint: DiscordEndpoint,
	}
}

// HandleDiscordLogin redirige l'utilisateur vers la page de consentement Discord
// C'est la première étape du processus OAuth : demander l'autorisation à l'utilisateur
func HandleDiscordLogin(w http.ResponseWriter, r *http.Request) {
	// Vérification que la configuration Discord a bien été initialisée
	if DiscordOauthConfig == nil {
		logMsg := "ERREUR : <discord.go> La communication discord n'a pas été initialisée."
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}
	// Génération de l'URL d'autorisation avec un token d'état pour la sécurité
	url := DiscordOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// HandleDiscordCallback gère la redirection de retour depuis Discord après autorisation
func HandleDiscordCallback(w http.ResponseWriter, r *http.Request) {
	// ÉTAPE 1 : Récupération du code d'autorisation depuis l'URL
	code := r.URL.Query().Get("code")
	if code == "" {
		logMsg := "ERREUR : <discord.go> Erreur dans la tentative de connexion, Discord n'a pas renvoyé de code d'autorisation."
		logs.AddLogsToDatabase(logMsg)
		utils.StatusBadRequest(w)
		return
	}

	// ÉTAPE 2 : Échange du code d'autorisation contre un token d'accès
	token, err := DiscordOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <discord.go> Erreur dans l'utilisation du code d'autorisation : ", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}

	// ÉTAPE 3 : Récupération des informations utilisateur via l'API Discord
	resp, err := http.Get("https://discord.com/api/v10/users/@me?access_token=" + token.AccessToken)
	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <discord.go> Impossible de récupérer les données de l'utilisateur : ", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}
	defer resp.Body.Close()

	// Décodage de la réponse JSON contenant les informations utilisateur
	var userInfo map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&userInfo)

	// ÉTAPE 4 : Extraction et validation des données essentielles
	// Récupération de l'ID Discord (identifiant unique de l'utilisateur chez Discord)
	discordID, ok := userInfo["id"].(string)
	if !ok {
		logMsg := "ERREUR : <discord.go> ID utilisateur Discord manquant"
		logs.AddLogsToDatabase(logMsg)
		return
	}

	// Récupération de l'email (obligatoire pour notre système)
	email, ok := userInfo["email"].(string)
	if !ok || email == "" {
		logMsg := "ERREUR : <discord.go> Email utilisateur Discord manquant/non autorisé"
		logs.AddLogsToDatabase(logMsg)
		return
	}

	// Récupération du nom d'utilisateur Discord
	username, ok := userInfo["username"].(string)
	if !ok {
		username = "DiscordUser" // Nom par défaut si non fourni
	}

	// ÉTAPE 5 : Recherche ou création de l'utilisateur dans la base de données locale
	userID, err := DiscordUser(discordID, email, username)
	if err != nil {
		logMsg := fmt.Sprint("Échec de la recherche/création de l'utilisateur : ", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}

	// ÉTAPE 6 : Création de la session utilisateur (cookie)
	err = authhandlers.InitSession(w, userID, "user", username)
	if err != nil {
		utils.InternalServError(w)
		return
	}

	// ÉTAPE 7 : Redirection vers la page d'accueil
	http.Redirect(w, r, "/", http.StatusFound)
}

// DiscordUser gère la logique de recherche ou de création d'un utilisateur dans la base de données locale
func DiscordUser(discordID, email, username string) (int, error) {
	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		return 0, err
	}
	defer db.Close()

	var userID int

	// CAS 1 : Recherche d'un utilisateur ayant déjà ce discord_id
	sqlQuery := `SELECT id FROM user WHERE discord_id = ?`
	row := db.QueryRow(sqlQuery, discordID)
	err = row.Scan(&userID)

	if err == nil {
		// L'utilisateur a été trouvé avec ce discord_id, on renvoie son ID
		return userID, nil
	} else if err != sql.ErrNoRows {
		// Erreur inattendue dans la base de données
		return 0, err
	}

	// CAS 2 et 3 : L'utilisateur n'a pas lié son compte Discord
	// On vérifie s'il n'a pas utilisé cette adresse email pour créer un compte
	if err == sql.ErrNoRows {
		// Recherche d'un utilisateur avec cette adresse email
		sqlQuery = `SELECT id FROM user WHERE email = ?`
		row = db.QueryRow(sqlQuery, email)
		err = row.Scan(&userID)

		switch err {
		// CAS 2 : L'utilisateur existe avec cet email → on associe son discord_id
		// Cela permet à l'utilisateur de se connecter via Discord à l'avenir
		case nil:
			sqlUpdate := `UPDATE user SET discord_id = ? WHERE id = ?`
			_, err = db.Exec(sqlUpdate, discordID, userID)
			if err != nil {
				return 0, err
			}
		// CAS 3 : Aucun utilisateur n'existe avec cette adresse mail ou ce discord_id
		// On crée un nouveau compte dans la base de données
		case sql.ErrNoRows:
			userID, err = CreateNewDiscordUser(discordID, email, username, db)
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

// CreateNewDiscordUser crée un nouvel utilisateur dans la base de données avec ses informations Discord
func CreateNewDiscordUser(discordID, email, discordName string, db *sql.DB) (int, error) {
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
	uniqueUsername := discordName
	for {
		var id int
		testedName := discordName
		if addon != 0 {
			// Construction du nom avec suffixe : nom_1, nom_2, etc.
			testedName = fmt.Sprintf("%s_%d", discordName, addon)
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
			uniqueUsername = testedName
			break
		}
	}

	// ÉTAPE 3 : Insertion du nouvel utilisateur dans la base de données
	// Note : la table 'user' a une nouvelle colonne 'discord_id'
	sqlUpdate := `INSERT INTO user(username, email, discord_id, role_id) VALUES(?, ?, ?, ?)`
	result, err := db.Exec(sqlUpdate, uniqueUsername, email, discordID, role)
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
