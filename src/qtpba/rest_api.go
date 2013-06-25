package qtpba

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func SetupRouter(router *mux.Router) {
	router.HandleFunc("/api/coordinates", serveCoordinates)
	router.HandleFunc("/api/tweet/{id}", serveTweet)
	router.HandleFunc("/api/tweetsByUser/{id}", serveTweetsByUser)
	router.HandleFunc("/api/checkins", serveCheckins)
	router.Handle("/", http.FileServer(http.Dir(GetFullPath("/static"))))
}

type JSONCoordinate struct {
	TweetId   int64
	Lat, Long float32
}

func serveTweets(resp http.ResponseWriter, req *http.Request, query string, params ...interface{}) {
	db := GetDB()
	rows, err := db.Query(query, params...)
	if err != nil {
		panic(err)
	}

	coordinates := make([]JSONCoordinate, 0)
	var coordinate JSONCoordinate
	for rows.Next() {
		rows.Scan(&coordinate.TweetId, &coordinate.Lat, &coordinate.Long)
		coordinates = append(coordinates, coordinate)
	}

	encoder := json.NewEncoder(resp)
	encoder.Encode(coordinates)
}

func serveCoordinates(resp http.ResponseWriter, req *http.Request) {
	serveTweets(resp, req, "SELECT id, latitude, longitude FROM tweet ORDER BY id DESC LIMIT 1000")
}

func serveCheckins(resp http.ResponseWriter, req *http.Request) {
	serveTweets(resp, req, "SELECT id, latitude, longitude FROM tweet WHERE message LIKE 'I''m at %' ORDER BY id DESC LIMIT 1000")
}

func serveTweetsByUser(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	if userId, err := strconv.ParseInt(vars["id"], 10, 64); err == nil { // Numeric ID
		serveTweets(resp, req, "SELECT id, latitude, longitude FROM tweet WHERE user_id=? ORDER BY id ASC", userId)
	} else { // String ID (user name)
		serveTweets(resp, req, "SELECT id, latitude, longitude FROM tweet WHERE user_id=(SELECT id FROM user WHERE screen_name=?) ORDER BY id ASC", id)
	}
}

func serveTweet(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	tweetId, _ := strconv.ParseInt(vars["id"], 10, 64)

	tweet := Tweet{}
	encoder := json.NewEncoder(resp)

	row := GetDB().QueryRow("SELECT user_id, message FROM tweet WHERE id=?", tweetId)
	if err := row.Scan(&tweet.User.Id, &tweet.Text); err != nil {
		if err == sql.ErrNoRows {
			encoder.Encode(nil)
		} else {
			encoder.Encode(err)
		}
		return
	}

	encoder.Encode(tweet)
}
