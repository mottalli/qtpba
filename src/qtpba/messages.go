package qtpba

import (
	"errors"
	"time"
)

func processMessage(message *ServerMessage) {
	var tweet *Tweet
	var err error

	if tweet, err = tweetFromJSON(message); err != nil {
		logger.Println(err)
		return
	}

	logger.Println(tweet)

	db.SaveTweet(tweet)
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

func tweetFromJSON(rawTweet *ServerMessage) (tweet *Tweet, err error) {
	tweet = new(Tweet)
	err = nil

	var ok bool

	if text, ok := (*rawTweet)["text"]; ok {
		if tweet.Text, ok = text.(string); !ok {
			err = errors.New("Field 'text' cannot be converted to a string")
			return
		}
	} else {
		err = errors.New("Field 'text' not found in server message")
		return
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

	return
}
