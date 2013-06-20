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

	db = &database{conn}
}

func GetDB() *database {
	return db
}

func (db *database) SaveTweet(tweet *Tweet) {
	stmt, _ := db.Prepare("INSERT INTO tweets(user, message, lat, long, timestamp_utc) VALUES(?, ?, ?, ?, ?)")
	stmt.Exec(tweet.User, tweet.Text, tweet.Coordinate.Lat, tweet.Coordinate.Long, tweet.Timestamp.Unix())
	stmt.Close()
}
