create table if not exists users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    tgid INTEGER UNIQUE NOT NULL,
    username text NOT NULL,
    rights text NOT NULL
);

create table if not exists restoraunts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    rest_name text NOT NULL,
    map_url text NOT NULL,
    reference_by INTEGER NOT NULL,
    votes INTEGER DEFAULT 0,
    FOREIGN KEY (reference_by) REFERENCES users(id)
 );

create table if not exists reviews (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    restoraunt_id INTEGER NOT NULL,
    category text NOT NULL,
    rate INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (restoraunt_id) REFERENCES restoraunts(id)
);

create table if not exists events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    visite_date DATE NOT NULL,
    rest_id INTEGER NOT NULL,
    FOREIGN KEY (rest_id) REFERENCES restoraunts(id)
);

create table if not exists possible_dates (
    possible_date DATE PRIMARY KEY,
    votes INTEGER DEFAULT 0
);

create table if not exists stage (
    current_stage INTEGER NOT NULL
);

create table if not exists actions {
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    action_type TEXT NOT NULL,
    subject_type TEXT DEFAULT NULL,
    subject_id INTEGER DEFAULT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
}