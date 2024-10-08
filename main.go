package main

import (
	"net/http"
	"source-query-test/async"
	"time"
)

func main() {
	go async.StartTimer()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		header.Set("Content-Type", "audio/mp3")
		header.Set("Access-Control-Allow-Origin", "*")

		w.WriteHeader(http.StatusOK)

		async.ResponseWriters[w] = true

		notify := r.Context().Done()

		for {
			select {
			case <-notify:
				delete(async.ResponseWriters, w)
				return
			default:
				time.Sleep(1 * time.Second)
			}
		}
	})

	http.ListenAndServe(":8080", nil)
}
