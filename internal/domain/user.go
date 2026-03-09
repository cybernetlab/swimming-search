package domain

import (
	"strings"
)

type User struct {
	Name      string
	IsAdmin   bool
	CentreIDs []uint
}

func (u User) Authorize(action string) bool {
	if u.Name == "" {
		return false
	}
	if u.IsAdmin {
		return true
	}
	if strings.Index(action, "user:") == 0 {
		return false
	}
	return true
}
