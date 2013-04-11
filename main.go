package main

import (
	_ "github.com/mattn/go-sqlite3"
    "fmt"
    "log"
    "os"
    "net/http"
)

var logger *log.Logger

func handleError(e error) {
	if e != nil {
		panic("Runtime error: " + e.Error())
	}
}

func processTweet(tweet *Tweet) {
    fmt.Println(tweet)
    saveTweet(tweet)
}

func listenForTweets() {
    tweets := make(chan(*Tweet))
    go runTweetStream(tweets)
    for {
        tweet := <-tweets
        processTweet(tweet)
    }
}

func main() {
	var err error

    logfile, err := os.OpenFile("./process.log", os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0666)
    handleError(err)
    logger = log.New(logfile, "", log.LstdFlags)
    logger.Println("Initializing app")

	err = initializeDB()
	handleError(err)
	defer db.Close()

    go listenForTweets()
    go runStatsDaemon()

    panic(http.ListenAndServe(":8080", http.FileServer(http.Dir("./"))))
}
