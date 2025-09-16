package utils

import (
	"database/sql"
	"log"

	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/utils/getdata"
)

func ChangeDisLikes(db *sql.DB, userID int, post models.Message) error {
	notliked, err := getdata.CheckIfLiked(db, post.MessageID, userID)
	if err != nil {
		log.Print(err)
		return err
	}

	var notdisliked bool
	notdisliked, err = getdata.CheckIfDisliked(db, post.MessageID, userID)
	if err != nil {
		log.Print(err)
		return err
	}

	// Stocke le nombre de likes et de dislikes dans des variables temporaires
	newlikes := post.Likes
	newdislikes := post.Dislikes
	if notliked && notdisliked {
		// Si le post n'était pas déjà liké ou disliké, rajouté un dislike
		newdislikes += 1
	} else if !notliked && notdisliked {
		// Si le post était liké, retire le like et ajoute un dislike
		newdislikes += 1
		newlikes -= 1
	} else if notliked && !notdisliked {
		// Si le post était disliké, retire le dislike
		newdislikes -= 1
	}

	// Vérifie les changements et met à jour les bases de données likes et dislikes
	if newdislikes > post.Dislikes {
		// Si un dislike a été ajouté, ajoute le post dans la base de données des dislikes
		if err := AddLikesAndDislikes(db, post.MessageID, userID, "dislikes"); err != nil {
			return err
		}
		if newlikes < post.Likes {
			// Si le post était liké avant d'être disliké, retire le post de la liste des likes
			if err := RemoveLikesAndDislikes(db, post.MessageID, userID, "likes"); err != nil {
				return err
			}

		}
	} else if newdislikes < post.Dislikes {
		// Si le post était déjà disliké, annule le dislike et le retire de la liste
		if err := RemoveLikesAndDislikes(db, post.MessageID, userID, "dislikes"); err != nil {
			return err
		}
	}

	// Met à jour la base de données pour le message disliké
	if err = UpdateLikesAndDislikes(db, post.MessageID, userID, newlikes, newdislikes, "dislikes"); err != nil {
		return err
	}

	// Met à jour la structure pour la renvoyer au handler
	post.Likes = newlikes
	post.Dislikes = newdislikes

	return nil
}
