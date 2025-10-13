package subhandlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/sessions"
	"github.com/Mathis-Pain/Forum/utils"
	"github.com/Mathis-Pain/Forum/utils/logs"
	"github.com/Mathis-Pain/Forum/utils/postactions"
)

// Gestion des likes sur les posts
func LikePostHandler(w http.ResponseWriter, r *http.Request) {
	// Récupère l'ID de l'utilisateur connecté et les infos du post liké
	userID, likedPost, err := getSessionAndPostInfo(r)
	if err != nil {
		utils.InternalServError(w)
		return
	}

	postactions.ChangeLikes(userID, likedPost)
	url := fmt.Sprintf("/topic/%d#%d", likedPost.TopicID, likedPost.MessageID)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

// Gestion des dislikes
func DislikePostHandler(w http.ResponseWriter, r *http.Request) {
	// Récupère l'ID de l'utilisateur connecté et les infos du post disliké
	userID, likedPost, err := getSessionAndPostInfo(r)
	if err != nil {
		utils.InternalServError(w)
		return
	}

	postactions.ChangeDisLikes(userID, likedPost)
	url := fmt.Sprintf("/topic/%d#%d", likedPost.TopicID, likedPost.MessageID)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

// Récupère les données de l'utilisateur et celles du post pour pouvoir mettre à jour les likes et dislikes
func getSessionAndPostInfo(r *http.Request) (int, models.Message, error) {
	// Récupère l'ID de l'utilisateur connecté
	cookie, _ := r.Cookie("session_id")
	session, err := sessions.GetSession(cookie.Value)
	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <likesdislikes.go> Erreur dans la récupération de session : ", err)
		logs.AddLogsToDatabase(logMsg)
		return 0, models.Message{}, err
	}
	userID := session.UserID

	postID, _ := strconv.Atoi(r.FormValue("postID"))
	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		logMsg := fmt.Sprintf("ERREUR : <likesdislikes.go> Erreur à l'ouverture de la base de données : %v\n", err)
		logs.AddLogsToDatabase(logMsg)
		return 0, models.Message{}, err
	}
	defer db.Close()

	post, err := postactions.GetMessageLikesAndDislikes(db, postID)
	if err != nil {
		logMsg := fmt.Sprint("ERREUR : <likesdislikes.go> Erreur dans la récupération des Likes/Dislikes :", err)
		logs.AddLogsToDatabase(logMsg)
		return userID, models.Message{}, err
	}

	return userID, post, nil
}
