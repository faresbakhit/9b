CREATE TABLE user(
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    username        TEXT UNIQUE NOT NULL,
    hashed_password BLOB NOT NULL,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    session_token   TEXT
);

CREATE TABLE post(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id    INTEGER,
    title      TEXT NOT NULL,
    content    TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY(user_id) REFERENCES user(id)
);
