package types

import "time"

type Discord_message struct {
	Id        string
	Published *time.Time
	Content   string
	Locale    string
	Title     string
	Image     string
}
