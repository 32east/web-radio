package handlers

import (
	"encoding/json"
	"net/http"
	"source-query-test/lib"
	"time"
)

func GetTrackInfo(w http.ResponseWriter, r *http.Request) {
	var currentTime = float64(time.Now().Sub(lib.CurrentMusic.StartTime).Milliseconds()) / 1000
	var musicInfo = lib.MusicInfo{
		Name:     lib.CurrentMusic.Name,
		Time:     currentTime,
		Duration: lib.CurrentMusic.Duration.Milliseconds() / 1000,
	}

	json.NewEncoder(w).Encode(musicInfo)
}

func GetListeners(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]int{
		"count": len(lib.Writers),
	})
}
