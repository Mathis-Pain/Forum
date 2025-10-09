-- Table des types
CREATE TABLE type (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `type` TEXT NOT NULL UNIQUE
);

-- Table des contenus de notifications
CREATE TABLE `notifications` (
    `id` INTEGER PRIMARY KEY, 
    `receiver_id` INTEGER NOT NULL, 
    `type` INTEGER REFERENCES `type`(`ID`), 
    `message` TEXT
);

-- Table des logs
CREATE TABLE `logs` (
    `id` INTEGER PRIMARY KEY, 
    `message` TEXT,
    `type` TEXT, 
    `date` TEXT DEFAULT (CURRENT_TIMESTAMP)
);
