package appwindow

import "strings"

func upString(s string) string {
	return strings.Title(strings.ToLower(strings.Replace(s, "_", " ", -1)))
}
