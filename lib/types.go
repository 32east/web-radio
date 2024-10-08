package lib

import "time"

type Music struct {
	Name            string
	StartTime       time.Time
	Duration        time.Duration
	Content         []byte
	LastEndPosition int
}

type MusicInfo struct {
	Name     string  `json:"name"`
	Time     float64 `json:"time"`
	Duration int64   `json:"duration"`
}
