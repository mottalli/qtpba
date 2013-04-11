package main

import (
	"time"
	"fmt"
    "net/http"
    "encoding/json"
    "os"
)

type RawTweet map[string]interface{}

type Coordinate struct {
	Lat, Long float64
}

type Tweet struct {
	Text, User string
	Coordinate Coordinate
	Timestamp  time.Time
}

func (tweet *Tweet) String() string {
	return fmt.Sprintf("[%v] @%v: \"%v\" %v", tweet.Timestamp, RED(tweet.User), tweet.Text, YELLOW(tweet.Coordinate))
}

func coordinateFromJSON(json interface{}) (coord Coordinate) {
	var coordinateMap map[string]interface{}
	var coordinatePair []interface{}
	var ok bool

	if coordinateMap, ok = json.(map[string]interface{}); !ok {
		return
	}

	if coordinatePair, ok = coordinateMap["coordinates"].([]interface{}); !ok {
		return
	}

	coord.Lat, _ = coordinatePair[1].(float64)
	coord.Long, _ = coordinatePair[0].(float64)
	return
}

func tweetFromJSON(rawTweet *RawTweet) (*Tweet, bool) {
	tweet := &Tweet{}
	var ok bool

	if text, ok := (*rawTweet)["text"]; ok {
		if tweet.Text, ok = text.(string); !ok {
			return nil, false
		}
	} else {
		return nil, false
	}

	// Me fijo si puedo sacar el usuario
	if _, ok = (*rawTweet)["user"]; ok {
		if userMap, ok := (*rawTweet)["user"].(map[string]interface{}); ok {
			if screenName, ok := userMap["screen_name"].(string); ok {
				tweet.User = screenName
			}
		}
	}

	if coordinate, ok := (*rawTweet)["coordinates"]; ok && coordinate != nil {
		tweet.Coordinate = coordinateFromJSON(coordinate)
	}

	if createdAt, ok := (*rawTweet)["created_at"]; ok {
		if timeString, ok := createdAt.(string); ok {
			tweet.Timestamp, _ = time.Parse("Mon Jan 2 15:04:05 -0700 2006", timeString)
		}
	}

	return tweet, true
}

func runTweetStream(stream chan<-(*Tweet)) {
    client := &http.Client{}
    url := "https://stream.twitter.com/1.1/statuses/filter.json?locations=-59,-35,-58,-34&track=buenos+aires,bsas,baires,bs.as.,bs.as"

    username, password := os.Getenv("QTPBA_USERNAME"), os.Getenv("QTPBA_PASSWORD")
    if username == "" || password == "" {
        logger.Fatalln("Either QTPBA_USERNAME or QTPBA_PASSWORD environment variables not set")
    }

    for {
        req, err := http.NewRequest("POST", url, nil)
        handleError(err)
        req.SetBasicAuth(username, password)

        var resp *http.Response

        if resp, err = client.Do(req); err != nil {
            logger.Println("Error requesting URL. Trying again in 10 seconds...")
            time.Sleep(10 * time.Second)
            continue
        }

        rawTweet := RawTweet{}
        jsonDecoder := json.NewDecoder(resp.Body)
        for {
            err = jsonDecoder.Decode(&rawTweet)
            if err != nil {
                fmt.Println("Error decoding JSON stream:", err, ". Retrying connection.")
                break
            }

            tweet, ok := tweetFromJSON(&rawTweet)
            if !ok {
                logger.Println("No se pudo decodificar el tweet crudo:", rawTweet)
                continue
            }

            stream <- tweet
        }
        resp.Body.Close()
    }
}
