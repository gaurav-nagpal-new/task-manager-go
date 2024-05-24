package model

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	Created    = "Created"
	InProgress = "In Progress"
	Done       = "Done"
)

func IsValidStatus(status string) bool {
	switch cases.Title(language.English, cases.Compact).String(status) {
	case Created:
		return true
	case InProgress:
		return true
	case Done:
		return true
	case "":
		return true
	default:
		return false
	}
}

func IsValidPriority(priority int) bool {
	switch priority {
	case -1:
		return true
	case 1:
		return true
	case 0:
		return true
	default:
		return false
	}
}
