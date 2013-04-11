package main

import (
	"database/sql"
)

var db *sql.DB

func initializeDB() (err error) {
	db, err = sql.Open("sqlite3", "./qtpba.db")
    return
}

func saveTweet(tweet *Tweet) {
	stmt, _ := db.Prepare("INSERT INTO tweets(user, message, lat, long, timestamp_utc) VALUES(?, ?, ?, ?, ?)")
	stmt.Exec(tweet.User, tweet.Text, tweet.Coordinate.Lat, tweet.Coordinate.Long, tweet.Timestamp.Unix())
	stmt.Close()
}
