package subhandlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/sessions"
	"github.com/Mathis-Pain/Forum/utils"
	"github.com/Mathis-Pain/Forum/utils/getdata"
)

// Gestion des likes sur les posts
func LikePostHandler(w http.ResponseWriter, r *http.Request) {
	// Récupère l'ID de l'utilisateur connecté et les infos du post liké
	userID, likedPost, err := getSessionAndPostInfo(r)
	if err != nil {
		utils.InternalServError(w)
		return
	}

	utils.ChangeLikes(userID, likedPost)

	url := fmt.Sprintf("/topic/%d", likedPost.MessageID)
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

	utils.ChangeDisLikes(userID, likedPost)
	url := fmt.Sprintf("/topic/%d", likedPost.MessageID)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func getSessionAndPostInfo(r *http.Request) (int, models.Message, error) {
	// Récupère l'ID de l'utilisateur connecté
	cookie, _ := r.Cookie("session_id")
	session, err := sessions.GetSession(cookie.Value)
	if err != nil {
		log.Print("<likesdislikes.go> Erreur dans la récupération de session : ", err)
		return 0, models.Message{}, err
	}
	userID := session.UserID

	postID, _ := strconv.Atoi(r.FormValue("postID"))
	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		log.Printf("<topichandler.go> Could not open database : %v\n", err)
		return 0, models.Message{}, err
	}
	defer db.Close()

	post, err := getdata.GetMessageLikesAndDislikes(db, postID)
	if err != nil {
		log.Print("<likesdislikes.go> Erreur dans la récupération des Likes/Dislikes :", err)
	}

	return userID, post, nil
}
