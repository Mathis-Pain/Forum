-- Table des types
CREATE TABLE role (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type TEXT NOT NULL UNIQUE
);

-- Table des contenus de notifications
CREATE TABLE `notifications` (
    `receiver_id` INTEGER PRIMARY KEY, 
    `type` INTEGER REFERENCES `type`(`ID`), 
    `message` TEXT
);