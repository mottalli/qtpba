CREATE TABLE user(
    id                  INTEGER PRIMARY KEY,
    screen_name         TEXT,
    name                TEXT,
    description         TEXT,
    followers_count     INTEGER,
    friends_count       INTEGER,
    language            TEXT,
    location            TEXT
);

CREATE TABLE tweet(
    id                  INTEGER PRIMARY KEY, 
    user_id             INTEGER,
    message             TEXT, 
    latitude            FLOAT, 
    longitude           FLOAT, 
    timestamp_utc       INT,
    FOREIGN KEY(user_id) REFERENCES user(id)
);

