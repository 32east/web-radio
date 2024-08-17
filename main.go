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

		notify := w.(http.CloseNotifier).CloseNotify()

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

	server := &http.Server{Addr: ":80"}
	server.SetKeepAlivesEnabled(false)
	server.ListenAndServe()
}
