package qtpba

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path"
)

const DBFILE = "qtpba.db"

type database struct {
	*sql.DB
}

var db *database

func initDB() {
	var err error
	var conn *sql.DB

	dbFile := path.Join(GetBaseDir(), "db/", DBFILE)

	if _, err = os.Stat(dbFile); os.IsNotExist(err) {
		logger.Fatal(err)
	}

	if conn, err = sql.Open("sqlite3", dbFile); err != nil {
		logger.Fatal(err)
	}

	logger.Println("Opened database", dbFile)

	db = &database{conn}
}

func GetDB() *database {
	return db
}

func (db *database) SaveUser(user TwitterUser) {
	db.Exec("DELETE FROM user WHERE id=?", user.Id)
	db.Exec("INSERT INTO user(id, screen_name, name, description, followers_count, friends_count, language, location) VALUES(?, ?, ?, ?, ?, ?, ?, ?)",
		user.Id, user.ScreenName, user.Name, user.Description, user.FollowersCount, user.FriendsCount, user.Language, user.Location)
}

func (db *database) SaveTweet(tweet *Tweet) {
	db.SaveUser(tweet.User)
	db.Exec("INSERT INTO tweet(user_id, message, latitude, longitude, timestamp_utc) VALUES(?, ?, ?, ?, ?)",
		tweet.User.Id, tweet.Text, tweet.Coordinate.Lat, tweet.Coordinate.Long, tweet.Timestamp.Unix())
}
