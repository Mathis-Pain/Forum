package admin

import (
	"database/sql"
	"log"
	"sort"

	"github.com/Mathis-Pain/Forum/models"
	"github.com/Mathis-Pain/Forum/utils/getdata"
)

func GetAllTopics(categories []models.Category, db *sql.DB) ([]models.Category, []models.Topic, error) {
	for i := 0; i < len(categories); i++ {
		var err error
		categories[i], err = getdata.GetCatDetails(db, categories[i].ID)
		categories[i].Topics, err = getdata.GetTopicList(db, categories[i].ID)
		if err != nil && err != sql.ErrNoRows {
			return categories, nil, err
		}
	}

	var topics []models.Topic

	for i := 0; i < len(categories); i++ {
		topicList, err := getdata.GetTopicList(db, categories[i].ID)
		if err != nil {
			log.Print("<displaydashboard.go> Erreur dans la récupération des sujets :", err)
			return categories, nil, err
		}

		topicList, err = GetCatName(categories[i], db, topicList)
		if err != nil {
			log.Print("<displaydashboard.go> Erreur dans la récupération des noms de catégorie :", err)
			return categories, nil, err
		}

		topics = append(topics, topicList...)
	}

	sort.Slice(topics, func(i, j int) bool {
		return topics[i].TopicID > topics[j].TopicID
	})

	return categories, topics, nil
}

func GetCatName(categ models.Category, db *sql.DB, topicList []models.Topic) ([]models.Topic, error) {
	for j := 0; j < len(topicList); j++ {
		var catname string
		sqlQuery := `SELECT name FROM category WHERE id = ?`
		rows, err := db.Query(sqlQuery, categ.ID)
		if err != nil {
			return topicList, err
		}

		for rows.Next() {
			if err := rows.Scan(&catname); err != nil {
				return topicList, err
			}
		}

		topicList[j].CatName = catname
	}

	return topicList, nil
}

func GetStats(topics []models.Topic) ([]models.LastPost, models.Stats, []models.User, error) {
	var stats models.Stats
	var users []models.User
	var err error

	users, stats.TotalUsers, err = GetAllUsers()
	if err != nil {
		log.Print("<displaydashboard.go, GetStats> Erreur dans la récupération des utilisateurs", err)
		return nil, models.Stats{}, nil, err
	}

	index := len(users) - 1
	if index < 0 {
		index = 0
	}
	stats.LastUser = users[index].Username

	if len(topics) != 0 {
		index = len(topics) - 1
		if index < 0 {
			index = 0
		}
		stats.LastTopic = topics[index].Name
	}
	stats.TotalTopics = len(topics)

	var lastMonthPosts []models.LastPost
	lastMonthPosts, stats.LastMonthPost, err = getdata.LastMonthPost()

	return lastMonthPosts, stats, users, nil
}

func GetAllUsers() ([]models.User, int, error) {
	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		log.Print("<displaydashboard.go, GetAllUsers> Erreur à l'ouverture de la base de données : ", err)
		return nil, 0, err
	}
	defer db.Close()

	var users []models.User
	var totalUsers int

	sqlQuery := `SELECT MAX(id) FROM user`
	err = db.QueryRow(sqlQuery).Scan(&totalUsers)
	if err != nil && err != sql.ErrNoRows {
		log.Print("<displaydashboard.go, GetAllUsers> Erreur dans la récupération du dernier ID utilisateur : ", err)
		return nil, 0, err
	}

	for i := 1; i <= totalUsers; i++ {
		user, err := getdata.GetUserInfoFromID(db, i)
		if err != nil {
			log.Print("<displaydashboard.go, GetAllUsers> Erreur dans la récupération des données utilisateurs : ", err)
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, totalUsers, nil
}
