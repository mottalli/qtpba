package qtpba

import (
	"fmt"
	"time"
)

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
