package getdata

import (
	"database/sql"

	"github.com/Mathis-Pain/Forum/models"
)

// Vérifie si l'utilisateur n'a pas liké le post
// Renvoie "true" si le post a été liké
func CheckIfLiked(db *sql.DB, postID, userID int) (bool, error) {
	// Vérifie si l'utilisateur n'a pas déjà liké ce post
	sqlQuery := `SELECT message_id FROM like WHERE user_id = ?`
	rows, err := db.Query(sqlQuery, userID)

	if err != nil {
		// Erreur dans la base de données
		return true, err
	}

	var likes models.Likes
	// Récupère la liste de tous les posts likés par l'utilisateur
	for rows.Next() {
		likedID := 0
		err := rows.Scan(&likedID)
		if err == sql.ErrNoRows {
			// S'il ne trouve ausun post dans la table like, renvoie false aussitôt (le post n'a jamais été liké par l'utilisateur)
			return false, nil
		} else if err != nil {
			// Erreur dans la base de données
			return true, err
		}
		likes.LikedPosts = append(likes.LikedPosts, likedID)
	}

	// Vérifie si le post actuel est présent dans la liste des posts likés par l'utilisateur
	for _, n := range likes.LikedPosts {
		if n == postID {
			return true, nil
		}
	}

	return false, nil

}

// Vérifie si l'utilisateur n'a pas disliké le post
// Renvoie "true" si le post a été dislike
func CheckIfDisliked(db *sql.DB, postID, userID int) (bool, error) {
	sqlQuery := `SELECT message_id FROM dislike WHERE user_id = ?`
	rows, err := db.Query(sqlQuery, userID)

	if err != nil {
		// Erreur dans la base de données
		return false, err
	}

	var dislikes models.Likes
	// Récupère la liste de tous les posts dislikés par l'utilisateur
	for rows.Next() {
		dislikedID := 0
		err := rows.Scan(&dislikedID)
		if err == sql.ErrNoRows {
			// S'il ne trouve ausun post dans la table dislike, renvoie false aussitôt (le post n'a jamais été disliké par l'utilisateur)
			return false, nil
		} else if err != nil {
			// Erreur dans la base de données
			return false, err
		}
		dislikes.LikedPosts = append(dislikes.LikedPosts, dislikedID)
	}

	// Vérifie si le post actuel est présent dans la liste des posts dislikés par l'utilisateur
	for _, n := range dislikes.LikedPosts {
		if n == postID {
			return true, nil
		}
	}

	return false, nil
}
