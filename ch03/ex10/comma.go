package ex10

import "bytes"

func Comma(s string) string {
	if len(s) <= 3 {
		return s
	}
	var buf bytes.Buffer
	for i := 0; i < len(s); i++ {
		if len(s[i:len(s)])%3 == 0 && buf.Len() > 0 {
			buf.WriteRune(',')
		}
		buf.WriteRune(rune(s[i]))
	}
	return buf.String()
}

func RecursiveComma(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}
	return RecursiveComma(s[:n-3]) + "," + s[n-3:]
}

func OptimizedComma(s string) string {
	if len(s) <= 3 {
		return s
	}
	var buf bytes.Buffer
	bufWritten := false
	n := len(s)
	for i := 0; i < n; i++ {
		if (n-i)%3 == 0 && bufWritten {
			buf.WriteRune(',')
		}
		buf.WriteRune(rune(s[i]))
		bufWritten = true
	}
	return buf.String()
}
