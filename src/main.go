package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"qtpba"
)

const (
	DATABASE_SUBDIR = "db"
	LOGS_SUBDIR     = "logs"
)

func main() {
	conn, err := qtpba.NewTwitterConnection()
	if err != nil {
		panic(err)
	}
	conn.StartListeningTweets()

	staticDir := qtpba.GetFullPath("static")
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		panic("The static files path " + staticDir + " does not exist")
	}

	router := mux.NewRouter()
	qtpba.SetupRouter(router)

	fmt.Println(qtpba.GetTweetStats().TopN(10))

	http.Handle("/", router)
	panic(http.ListenAndServe(":8000", nil))
}
