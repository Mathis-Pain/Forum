package getdata

import (
	"database/sql"

	"github.com/Mathis-Pain/Forum/models"
)

func CheckIfLiked(db *sql.DB, postID, userID int) (bool, error) {
	// Vérifie si l'utilisateur n'a pas déjà liké ce post
	sqlQuery := `SELECT post_id FROM likes WHERE user_id = ?`
	rows, err := db.Query(sqlQuery, userID)

	if err != nil {
		return false, err
	}

	var likes models.Likes
	// Récupère la liste de tous les posts likés par l'utilisateur
	for rows.Next() {
		likedID := 0
		err := rows.Scan(&likedID)
		if err == sql.ErrNoRows {
			return true, nil
		} else if err != nil {
			return false, err
		}
		likes.LikedPosts = append(likes.LikedPosts, likedID)
	}

	// Vérifie si le post actuel n'est pas déjà présent dans la liste
	for _, n := range likes.LikedPosts {
		if n == postID {
			return false, nil
		}
	}

	return true, nil

}

// Vérifie si l'utilisateur n'a pas disliké le post
func CheckIfDisliked(db *sql.DB, postID, userID int) (bool, error) {
	sqlQuery := `SELECT post_id FROM dislikes WHERE user_id = ?`
	rows, err := db.Query(sqlQuery, userID)

	if err != nil {
		return false, err
	}

	var dislikes models.Likes
	// Récupère la liste de tous les posts dislikés par l'utilisateur
	for rows.Next() {
		dislikedID := 0
		err := rows.Scan(&dislikedID)
		if err == sql.ErrNoRows {
			return true, nil
		} else if err != nil {
			return false, err
		}
		dislikes.LikedPosts = append(dislikes.LikedPosts, dislikedID)
	}

	// Vérifie si le post actuel n'est pas déjà présent dans la liste
	for _, n := range dislikes.LikedPosts {
		if n == postID {
			return false, nil
		}
	}

	return true, nil
}
