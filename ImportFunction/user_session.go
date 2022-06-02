package ImportFunction

import (
	"time"
)

var sessions = map[string]Session{}

// each session contains the username of the user and the time at which it expires
type Session struct {
	Username string
	Expiry   time.Time
}

// we'll use this method later to determine if the session has expired
func (s Session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}
