CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    github_token TEXT
);

CREATE TABLE IF NOT EXISTS projects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    userId INTEGER,
    projectName TEXT NOT NULL,
    projectDescription TEXT,
    FOREIGN KEY (userId) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS commits (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    projectId INTEGER NOT NULL,
    message TEXT NOT NULL,
    rating INTEGER,
    date TEXT NOT NULL DEFAULT (datetime('now')),
    FOREIGN KEY (projectId) REFERENCES projects(id) ON DELETE CASCADE
);
