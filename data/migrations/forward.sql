create table if not exists users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    tgid INTEGER UNIQUE,
    username text,
    rights text
)