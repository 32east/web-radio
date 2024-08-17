package async

import "time"

type Music struct {
	Name            string
	StartTime       time.Time
	Duration        time.Duration
	Content         []byte
	CurrentContent  []byte
	LastEndPosition int
}
