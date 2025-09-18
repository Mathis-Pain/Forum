package utils

import (
	"database/sql"
	"log"

	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/utils/getdata"
)

// Fonction pour mettre à jour le nombre de likes
func ChangeLikes(userID int, post models.Message) error {
	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		log.Printf("<topichandler.go> Could not open database : %v\n", err)
		return err
	}
	defer db.Close()

	// Vérifie si le post a déjà été liké par l'utilisateur connecté
	liked, err := getdata.CheckIfLiked(db, post.MessageID, userID)
	if err != nil {
		log.Print(err)
		return err
	}

	var disliked bool
	// Vérifie ensuite si le post a déjà été disliké par l'utilisateur
	disliked, err = getdata.CheckIfDisliked(db, post.MessageID, userID)
	if err != nil {
		log.Print(err)
		return err
	}

	// Stocke le nombre de likes et de dislikes dans des variables temporaires pour future comparaison
	newlikes := post.Likes
	newdislikes := post.Dislikes
	if !liked && !disliked {
		// Si le post n'était pas déjà liké ou disliké, rajoute un like
		newlikes += 1
	} else if !liked && disliked {
		// Si le post était disliké, retire le dislike et ajoute un like
		newdislikes -= 1
		newlikes += 1
	} else if liked && !disliked {
		// Si le post était liké, retire le like
		newlikes -= 1
	}

	// Vérifie les changements et met à jour les bases de données likes et dislikes
	if newlikes > post.Likes {
		// Si un like a été ajouté, ajoute le post dans la base de données des likes
		if err := AddLikesAndDislikes(db, post.MessageID, userID, "likes"); err != nil {
			return err
		}
		if newdislikes < post.Dislikes {
			// Si le post était liké avant d'être disliké, retire le post de la liste des dislikes
			if err := RemoveLikesAndDislikes(db, post.MessageID, userID, "dislikes"); err != nil {
				return err
			}

		}
	} else if newlikes < post.Likes {
		// Si le post était déjà liké, annule le like et le retire de la liste
		if err := RemoveLikesAndDislikes(db, post.MessageID, userID, "likes"); err != nil {
			return err
		}
	}

	// Met à jour la base de données pour le message liké
	if err = UpdateLikesAndDislikes(db, post.MessageID, userID, newlikes, newdislikes, "likes"); err != nil {
		return err
	}

	// Met à jour la structure pour la renvoyer au handler
	post.Likes = newlikes
	post.Dislikes = newdislikes

	return nil
}
