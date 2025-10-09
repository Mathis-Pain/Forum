package authhandlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/Mathis-Pain/Forum/handlers/subhandlers"
	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/utils"
	"github.com/Mathis-Pain/Forum/utils/logs"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// Si funcMap non declaré avant Funcs(funcMap) est not found
var funcMap2 = template.FuncMap{
	"toUpper": func(s string) string {
		return strings.ToUpper(s)
	},
}
var registrationHtml = template.Must(template.New("registration.html").Funcs(funcMap2).ParseFiles("templates/registration.html", "templates/login.html", "templates/header.html", "templates/initpage.html"))

func SignUpSubmitHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <cathandler.go> Erreur à l'ouverture de la base de données : %v\n", err)
		logs.AddLogsToDatabase(logMsg)
		return
	}
	defer db.Close()

	_, categories, _, err := subhandlers.BuildHeader(r, w, db)
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <cathandler.go> Erreur dans la construction du header : %v\n", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}

	if r.Method != http.MethodPost {
		data := struct {
			PageName     string
			Categories   []models.Category
			LoginErr     string
			CurrentUser  models.UserLoggedIn
			RegisterData models.RegisterDataError
			UserInfo     models.User
		}{
			PageName:     "Inscription",
			Categories:   categories,
			LoginErr:     "",
			CurrentUser:  models.UserLoggedIn{},
			RegisterData: models.RegisterDataError{},
			UserInfo:     models.User{},
		}
		// GET : afficher le formulaire vide
		if err := registrationHtml.Execute(w, data); err != nil {
			logMsg := fmt.Sprint("Erreur dans l'affichage de la page d'inscription :", err)
			logs.AddLogsToDatabase(logMsg)
			utils.InternalServError(w)
		}
		return
	}

	// --- Récupération des valeurs ---
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	passwordConfirm := r.FormValue("confirmpassword")
	profilPic := r.FormValue("userimage")

	if profilPic == "" {
		profilPic = "/static/noprofilpic.png"
	}

	// --- Struct pour stocker les erreurs ---

	formData := models.RegisterDataError{
		NameError:  utils.ValidName(username),
		EmailError: utils.ValidEmail(email),
		PassError:  utils.ValidPasswd(password, passwordConfirm),
	}

	userInfo := models.User{
		Username:  username,
		Email:     email,
		ProfilPic: profilPic,
	}

	data := struct {
		PageName     string
		Categories   []models.Category
		LoginErr     string
		CurrentUser  models.UserLoggedIn
		RegisterData models.RegisterDataError
		UserInfo     models.User
	}{
		PageName:     "Inscription",
		Categories:   categories,
		LoginErr:     "",
		CurrentUser:  models.UserLoggedIn{},
		RegisterData: formData,
		UserInfo:     userInfo,
	}

	// Si une erreur existe, renvoyer le formulaire avec messages
	if formData.NameError != "" || formData.EmailError != "" || formData.PassError != "" {
		w.WriteHeader(http.StatusBadRequest)
		registrationHtml.Execute(w, data)
		return
	}

	// --- Ouverture de la DB ---

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		utils.InternalServError(w)
		return
	}
	// ---- Vérifie si c'est le premier utilisateur ---
	var count int
	role := 3
	err = db.QueryRow("SELECT COUNT(*) FROM user").Scan(&count)

	if err != nil && err != sql.ErrNoRows {
		logMsg := fmt.Sprintf("ERREUR : Impossible de compter les utilisateurs existants : %v", err)
		logs.AddLogsToDatabase(logMsg)
		utils.InternalServError(w)
		return
	}

	if count == 0 {
		role = 1
	}

	// --- Insertion dans la DB ---
	_, err = db.Exec("INSERT INTO user(username, email, password, role_id) VALUES(?, ?, ?, ?)", username, email, hashedPassword, role)
	if err != nil {
		// Vérification UNIQUE (nom ou email déjà utilisé)
		if err.Error() == "UNIQUE constraint failed: user.username" {
			formData.NameError = "Ce nom d'utilisateur est déjà pris"
			w.WriteHeader(http.StatusBadRequest)
			registrationHtml.Execute(w, data)
			return
		} else if err.Error() == "UNIQUE constraint failed: user.email" {
			formData.EmailError = "Cette adresse email est déjà utilisée"
			w.WriteHeader(http.StatusBadRequest)
			registrationHtml.Execute(w, data)
			return
		}
		// Toute autre erreur
		utils.InternalServError(w)
		return
	}

	// --- Succès : redirection vers la page d'accueil ---
	logMsg := fmt.Sprintln("USER : Un nouvel utilisateur s'est inscrit : ", username)
	logs.AddLogsToDatabase(logMsg)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
