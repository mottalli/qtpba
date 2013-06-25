package qtpba

import (
	"fmt"
	"time"
)

type Coordinate struct {
	Lat, Long float64
}

type TwitterUser struct {
	Id             int64
	ScreenName     string
	Name           string
	Description    string
	FollowersCount int
	FriendsCount   int
	Language       string
	Location       string
}

type Tweet struct {
	Text       string
	User       TwitterUser
	Coordinate Coordinate
	Timestamp  time.Time
}

func (tweet *Tweet) String() string {
	return fmt.Sprintf("[%v] @%v: \"%v\" %v", tweet.Timestamp, RED(tweet.User.ScreenName), tweet.Text, YELLOW(tweet.Coordinate))
}
