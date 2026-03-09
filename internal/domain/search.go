package domain

import "context"

type Search struct {
	NameQuery string
	Days      []uint8
	CentreIDs []uint
	UserName  string
	ChatID    int64
	Cancel    context.CancelFunc `json:"-"`
}
