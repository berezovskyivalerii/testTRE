package model

import "testing"

func TestUser_IsBiz(t *testing.T) {
	cases := []struct {
		email string
		want  bool
	}{
		{"foo@bar.biz", true},
		{"FOO@BAR.BIZ", true},
		{"user@example.com", false},
	}
	for _, c := range cases {
		u := User{Email: c.email}
		if got := u.IsBiz(); got != c.want {
			t.Errorf("IsBiz(%q) = %v; want %v", c.email, got, c.want)
		}
	}
}