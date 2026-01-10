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
    closest_metro text NOT NULL,
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
    visit_date DATE DEFAULT NULL,
    rest_id INTEGER DEFAULT NULL,
    FOREIGN KEY (rest_id) REFERENCES restoraunts(id)
);

create table if not exists dates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    candidate DATE NOT NULL,
    event_id INTEGER NOT NULL,
    FOREIGN KEY (event_id) REFERENCES events(id)
);

create table if not exists votes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    rest_id INTEGER DEFAULT NULL,
    date_id INTEGER DEFAULT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (rest_id) REFERENCES restoraunts(id),
    FOREIGN KEY (date_id) REFERENCES possible_dates(id)
);

create table if not exists next_task (
    id INTEGER UNIQUE,
    dt DATE PRIMARY KEY
);

create table if not exists stage (
    current_stage INTEGER NOT NULL
);

create table if not exists actions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    dt DATETIME NOT NULL,
    user_id INTEGER NOT NULL,
    action_type TEXT NOT NULL,
    subject_type TEXT DEFAULT NULL,
    subject_id INTEGER DEFAULT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);