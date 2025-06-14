package model

import "strings"

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (u User) IsBiz() bool {
	return strings.HasSuffix(strings.ToLower(u.Email), ".biz")
}