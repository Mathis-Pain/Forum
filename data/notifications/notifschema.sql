-- Table des types
CREATE TABLE type (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `type` TEXT NOT NULL UNIQUE
);

-- Table des contenus de notifications
CREATE TABLE notifications (
    "ID" INTEGER PRIMARY KEY, 
    `receiver_id` INTEGER NOT NULL DEFAULT 1, 
    `type` INTEGER REFERENCES `type`(`ID`), 
    `message` TEXT, read INT DEFAULT 0, 
    "post_id" INT
);

-- Table des logs
CREATE TABLE logs (
    `ID` INTEGER PRIMARY KEY AUTOINCREMENT, 
    `message` TEXT, `type` TEXT, 
    `date` TEXT DEFAULT (CURRENT_TIMESTAMP), 
    handled INT DEFAULT 0, 
    sender INT
    );
