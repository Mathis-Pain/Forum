package utils

import (
	"database/sql"
	"log"

	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/utils/getdata"
)

func ChangeDisLikes(userID int, post models.Message) error {
	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		log.Printf("<topichandler.go> Could not open database : %v\n", err)
		return err
	}
	defer db.Close()

	liked, err := getdata.CheckIfLiked(db, post.MessageID, userID)
	if err != nil {
		log.Print(err)
		return err
	}

	var disliked bool
	disliked, err = getdata.CheckIfDisliked(db, post.MessageID, userID)
	if err != nil {
		log.Print(err)
		return err
	}

	// Stocke le nombre de likes et de dislikes dans des variables temporaires
	newlikes := post.Likes
	newdislikes := post.Dislikes
	if !liked && !disliked {
		// Si le post n'était pas déjà liké ou disliké, rajouté un dislike
		newdislikes += 1
	} else if liked && !disliked {
		// Si le post était liké, retire le like et ajoute un dislike
		newdislikes += 1
		newlikes -= 1
	} else if !liked && disliked {
		// Si le post était disliké, retire le dislike
		newdislikes -= 1
	}

	// Vérifie les changements et met à jour les bases de données likes et dislikes
	if newdislikes > post.Dislikes {
		// Si un dislike a été ajouté, ajoute le post dans la base de données des dislikes
		if err := AddLikesAndDislikes(db, post.MessageID, userID, "dislike"); err != nil {
			return err
		}
		if newlikes < post.Likes {
			// Si le post était liké avant d'être disliké, retire le post de la liste des likes
			if err := RemoveLikesAndDislikes(db, post.MessageID, userID, "like"); err != nil {
				return err
			}

		}
	} else if newdislikes < post.Dislikes {
		// Si le post était déjà disliké, annule le dislike et le retire de la liste
		if err := RemoveLikesAndDislikes(db, post.MessageID, userID, "dislike"); err != nil {
			return err
		}
	}

	// Met à jour la base de données pour le message disliké
	if err = UpdateLikesAndDislikes(db, post.MessageID, userID, newlikes, newdislikes, "dislike"); err != nil {
		return err
	}

	// Met à jour la structure pour la renvoyer au handler
	post.Likes = newlikes
	post.Dislikes = newdislikes

	return nil
}
