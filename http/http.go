package http

import (
	"net/http"
	"source-query-test/http/handlers"
	"source-query-test/lib"
	"source-query-test/middleware"
)

func SendImmediately(w *http.ResponseWriter) {
	var musicEndPos = lib.CurrentMusic.LastEndPosition
	var remap3Seconds = 3000 / float64(lib.CurrentMusic.Duration.Milliseconds())
	var minusSecond = float64(len(lib.CurrentMusic.Content)) * remap3Seconds
	var startPos = int64(float64(musicEndPos) - minusSecond)

	if startPos < 0 {
		startPos = 0
	}

	var startContent = lib.CurrentMusic.Content[startPos:musicEndPos]
	(*w).Write(startContent)
}

func Handle() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var header = w.Header()
		header.Set("Content-Type", "audio/mp3")
		header.Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)

		lib.Mutex.Lock()
		lib.Writers[&w] = true
		lib.Mutex.Unlock()

		SendImmediately(&w)
		<-r.Context().Done()

		lib.Mutex.Lock()
		delete(lib.Writers, &w)
		lib.Mutex.Unlock()
	})

	middleware.API("/api/v1/track-info", handlers.GetTrackInfo)
	middleware.API("/api/v1/listeners", handlers.GetListeners)

	http.ListenAndServe(":8080", nil)
}
