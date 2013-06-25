package qtpba

import (
	"errors"
	"strconv"
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

func userFromJSON(userMap map[string]interface{}) (user TwitterUser) {
	// Convert the ID from text to int64 (direct casting to int64 doesn't work)
	idStr := userMap["id_str"].(string)
	user.Id, _ = strconv.ParseInt(idStr, 10, 64)

	user.Language = userMap["lang"].(string)
	user.Name = userMap["name"].(string)
	user.ScreenName = userMap["screen_name"].(string)

	if description, ok := userMap["description"]; ok && description != nil {
		user.Description = description.(string)
	}

	if followersCount, ok := userMap["followers_count"]; ok && followersCount != nil {
		user.FollowersCount = int(followersCount.(float64))
	}

	if friendsCount, ok := userMap["friends_count"]; ok && friendsCount != nil {
		user.FriendsCount = int(friendsCount.(float64))
	}

	if location, ok := userMap["location"]; ok && location != nil {
		user.Location = location.(string)
	}

	return
}

func tweetFromJSON(rawTweet *ServerMessage) (tweet *Tweet, err error) {
	tweet = new(Tweet)
	err = nil

	var ok bool

	if text, ok := (*rawTweet)["text"]; ok {
		tweet.Text = text.(string)
	} else {
		err = errors.New("Field 'text' not found in server message (invalid tweet JSON)")
		return
	}

	// Me fijo si puedo sacar el usuario
	if _, ok = (*rawTweet)["user"]; ok {
		userMap := (*rawTweet)["user"].(map[string]interface{})
		tweet.User = userFromJSON(userMap)
	}

	if coordinate, ok := (*rawTweet)["coordinates"]; ok && coordinate != nil {
		tweet.Coordinate = coordinateFromJSON(coordinate)
	}

	if createdAt, ok := (*rawTweet)["created_at"]; ok {
		tweet.Timestamp, _ = time.Parse(time.RubyDate, createdAt.(string))
	}

	return
}
