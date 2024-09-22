CREATE TABLE posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    image TEXT,
    gif TEXT,
    privacy TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);