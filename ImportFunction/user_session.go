package ImportFunction

import (
	"net/http"
	"time"
)

type Cookie struct {
	Name       string
	Value      string
	Path       string
	Domain     string
	Expires    time.Time
	RawExpires string
	MaxAge     int
	Secure     bool
	HttpOnly   bool
	Raw        string
	Unparsed   []string
}

func CookieGenerator(username string) http.Cookie {
	expiration := time.Now().Add(time.Minute * 15)
	cookie := http.Cookie{Name: username, Value: username + "_cookie", Expires: expiration}
	return cookie
}
