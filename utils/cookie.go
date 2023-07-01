package utils

import "strings"

func ConvertStrCookieToDict(cookieStr string) map[string]string {
	cookieDict := make(map[string]string)

	if cookieStr == "" {
		return cookieDict
	}

	cookies := strings.Split(cookieStr, ";")
	for _, cookie := range cookies {
		cookie = strings.TrimSpace(cookie)
		if cookie == "" {
			continue
		}
		cookieArr := strings.Split(cookie, "=")
		cookieValue := cookieArr[1]
		if len(cookieValue) > 1 {
			cookieDict[cookieArr[0]] = cookieValue
		}
	}

	return cookieDict
}
