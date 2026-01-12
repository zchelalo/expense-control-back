package domain

import "strings"

type Email struct{ value string }

func NewEmail(v string) (Email, error) {
	v = strings.TrimSpace(strings.ToLower(v))
	if v == "" || len(v) > 254 || !looksLikeEmail(v) {
		return Email{}, ErrInvalidEmail
	}
	return Email{value: v}, nil
}

func (e Email) String() string { return e.value }

func looksLikeEmail(s string) bool {
	at := strings.IndexByte(s, '@')
	if at <= 0 {
		return false
	}
	if strings.Count(s, "@") != 1 {
		return false
	}
	dot := strings.LastIndexByte(s, '.')
	return dot > at+1 && dot < len(s)-1
}