package appwindow

import "strings"

func upString(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if i < 1 {
			if c >= 'a' && c <= 'z' {
				c -= 'a' - 'A'
			}
		}
		b.WriteByte(c)
	}
	return b.String()
}
