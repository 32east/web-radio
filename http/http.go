package http

import (
	"net/http"
	"source-query-test/lib"
)

func Handle() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var header = w.Header()
		header.Set("Content-Type", "audio/mp3")
		header.Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)

		lib.Writers[w] = true

		var notify = r.Context().Done()

		for {
			<-notify
			delete(lib.Writers, w)
		}
	})

	http.ListenAndServe(":8080", nil)
}
