package main

import (
	"log"
	"os"
	"os/signal"
	"source-query-test/http"
	"source-query-test/lib"
	"syscall"
)

func main() {
	var chn = make(chan os.Signal)
	signal.Notify(chn, syscall.SIGINT, syscall.SIGTERM)

	go lib.StartTimer()
	go http.Handle()

	<-chn
	log.Println("Радио отключено.")
}
