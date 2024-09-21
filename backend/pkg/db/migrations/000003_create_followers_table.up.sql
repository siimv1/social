CREATE TABLE followers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    follower_id INTEGER NOT NULL,
    followed_id INTEGER NOT NULL,
    status TEXT DEFAULT 'accepted', -- Jälgimise olek: accepted, pending, rejected
    FOREIGN KEY (follower_id) REFERENCES users(id),
    FOREIGN KEY (followed_id) REFERENCES users(id),
    UNIQUE(follower_id, followed_id) -- Tagab, et sama kasutaja ei saa sama kasutajat uuesti järgida
);
