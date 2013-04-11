CREATE TABLE IF NOT EXISTS tweets(
    id INTEGER PRIMARY KEY, 
    user TEXT, 
    message TEXT, 
    lat FLOAT, 
    long FLOAT, 
    timestamp_utc INT
);
