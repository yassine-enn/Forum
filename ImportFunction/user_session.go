package ImportFunction

import (
	"time"
)

// type Cookie struct {
// 	Name       string
// 	Value      string
// 	Path       string
// 	Domain     string
// 	Expires    time.Time
// 	RawExpires string
// 	MaxAge     int
// 	Secure     bool
// 	HttpOnly   bool
// 	Raw        string
// 	Unparsed   []string																																																																																																																																																																																																																																																																																																																																																																																																																								²²²²²²²²²²²
// }

// func CookieGenerator(username string) http.Cookie {
// 	expiration := time.Now().Add(time.Minute * 15)
// 	cookie := http.Cookie{Name: username, Value: username + "_cookie", Expires: expiration}
// 	return cookie
// }

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
