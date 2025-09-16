package utils

import (
	"database/sql"
	"log"

	"github.com/Mathis-Pain/Forum/models"
)

func ChangeLikes(db *sql.DB, postID, userID int, post models.Message) (models.Message, error) {
	// Vérifie si le post actuel a déjà été liké par l'utilisateur connecté
	notliked, err := CheckIfLiked(db, postID, userID)
	if err != nil {
		log.Print(err)
		return models.Message{}, err
	}

	// Vérifie si le post actuel a été disliké par l'utilisateur connecté
	var notdisliked bool
	notdisliked, err = CheckIfDisliked(db, postID, userID)
	if err != nil {
		log.Print(err)
		return models.Message{}, err
	}

	if notliked && notdisliked {
		// Si le post n'a été ni liké ni disliké, rajoute le like
		post.Likes += 1
	} else if notliked && !notdisliked {
		// Si le post était disliké, retire le dislike et rajoute le like
		post.Dislikes -= 1
		post.Likes += 1
	} else if !notliked {
		// Si le post est déjà liké, le compte de likes ne bouge pas
		log.Printf("<countlikes.go> L'utilisateur %d a déjà aimé ce post.", post.Author.ID)
		return post, nil
	}

	// Met à jour la base de données
	if UpdateLikesAndDislikes(db, postID, userID, post.Likes, post.Dislikes, "likes") != nil {
		return post, err
	}

	// Renvoie la structure mise à jour au serveur
	return post, nil
}
