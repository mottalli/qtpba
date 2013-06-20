package main

import (
	"net/http"
	"os"
	"qtpba"
)

const (
	DATABASE_SUBDIR = "db"
	LOGS_SUBDIR     = "logs"
)

func main() {
	go qtpba.StartListeningTweets()

	staticDir := qtpba.GetFullPath("static")
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		panic("The static files path " + staticDir + " does not exist")
	}

	//fmt.Println(qtpba.GetTweetStats().TopN(5))

	panic(http.ListenAndServe(":8000", http.FileServer(http.Dir(staticDir))))
}
