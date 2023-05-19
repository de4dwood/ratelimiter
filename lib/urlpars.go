package lib

import (
	"net/url"
	"strings"
)

func UrlParser(s string) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}
	if u.IsAbs() {
		host := u.Hostname()
		uri := strings.TrimSuffix(u.Path, "/")
		return host + uri, nil
	} else {
		return s, nil
	}
}
