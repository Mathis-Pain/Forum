
CREATE TABLE user (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT NOT NULL UNIQUE,
	email TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL
	profilpic TEXT DEFAULT './static/noprofilpic.png'
	role TEXT DEFAULT 
	 role_id INTEGER NOT NULL DEFAULT 1,
    FOREIGN KEY (role_id) REFERENCES roles(id)
);
CREATE TABLE category (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	description TEXT
);
CREATE TABLE topic (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	create_at TEXT DEFAULT (CURRENT_TIMESTAMP) NOT NULL,
	user_id INTEGER NOT NULL,
	category_id INTEGER NOT NULL,
	CONSTRAINT topic_category_FK FOREIGN KEY (category_id) REFERENCES category(id),
	CONSTRAINT topic_users_FK FOREIGN KEY (user_id) REFERENCES user(id)
);
CREATE TABLE message (
	id INTEGER NOT NULL,
	topic_id INTEGER NOT NULL,
	CONTENT TEXT NOT NULL,
	created_at TEXT DEFAULT (CURRENT_TIMESTAMP),
	likes INTEGER,
	dislikes INTEGER,
	user_id INTEGER NOT NULL,
	CONSTRAINT message_topic_FK FOREIGN KEY (topic_id) REFERENCES topic(id),
	CONSTRAINT message_users_FK FOREIGN KEY (user_id) REFERENCES user(id)
	);

CREATE TABLE roles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);
