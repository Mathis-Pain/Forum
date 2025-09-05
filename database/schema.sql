
CREATE TABLE user (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT DEFAULT (CURRENT_TIMESTAMP) NOT NULL,
	email TEXT NOT NULL,
	password TEXT NOT NULL
)
CREATE TABLE category (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	description TEXT
)
CREATE TABLE topic (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	create_at TEXT DEFAULT (CURRENT_TIMESTAMP) NOT NULL,
	user_id INTEGER NOT NULL,
	category_id INTEGER NOT NULL,
	CONSTRAINT topic_category_FK FOREIGN KEY (category_id) REFERENCES category(id),
	CONSTRAINT topic_users_FK FOREIGN KEY (user_id) REFERENCES "user"(id)
)
CREATE TABLE message (
	id INTEGER NOT NULL,
	topic_id INTEGER NOT NULL,
	"content" TEXT NOT NULL,
	created_at TEXT DEFAULT (CURRENT_TIMESTAMP),
	likes INTEGER,
	dislikes INTEGER,
	user_id INTEGER NOT NULL,
	CONSTRAINT message_topic_FK FOREIGN KEY (topic_id) REFERENCES topic(id),
	CONSTRAINT message_users_FK FOREIGN KEY (user_id) REFERENCES "user"(id)
)-- database: /Users/mathispain/Desktop/Zone01/forum/database/schema.db

SELECT * FROM "table-name";
