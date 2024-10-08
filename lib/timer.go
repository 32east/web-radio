package lib

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var CurrentMusic Music
var Writers = make(map[*http.ResponseWriter]bool)
var Mutex = sync.Mutex{}

func OnTimerTick() {
	var musicIsEnded = time.Now().After(CurrentMusic.StartTime.Add(CurrentMusic.Duration))

	if CurrentMusic.Name == "" || musicIsEnded {
		var files, err = os.ReadDir("music")

		if err != nil {
			log.Fatal("Error opening folder: ", err)
			return
		}

		var filesLength = len(files)
		var name = ""
		var fileName = ""

		for i := 0; i < 5; i++ {
			randomFile := files[rand.Intn(filesLength)]
			name = randomFile.Name()
			fileName = "music/" + name

			if name != CurrentMusic.Name {
				break
			}
		}

		var f, fileErr = os.Open(fileName)

		if fileErr != nil {
			log.Println("Error opening file: ", err)
			return
		}

		defer f.Close()
		var content, contentErr = GetFileContent(f)

		if contentErr != nil {
			log.Println("Error reading file: ", err)
			return
		}

		f.Seek(0, 0)

		duration := GetDuration(f)

		if duration == 0 {
			log.Println("Не удалось получись длительность трека: " + name)
			return
		}

		name = strings.TrimSuffix(name, filepath.Ext(name))
		CurrentMusic = Music{name, time.Now(), duration, content, 0}

		log.Println("Музыка изменена: ", name)
	} else {
		var musicContent = CurrentMusic.Content
		var contentLength = len(musicContent)
		var startPosition = CurrentMusic.LastEndPosition

		if startPosition >= contentLength {
			return
		}

		var currentTimePlusSecond = time.Now().Sub(CurrentMusic.StartTime).Milliseconds()
		var remapPlusSecond = float64(currentTimePlusSecond) / float64(CurrentMusic.Duration.Milliseconds())
		var endPosition = int(float64(contentLength) * remapPlusSecond)

		if endPosition > contentLength {
			endPosition = contentLength
		}

		CurrentMusic.LastEndPosition = endPosition

		var startContent = CurrentMusic.Content[startPosition:endPosition]

		Mutex.Lock()
		for ResponseWriter := range Writers {
			(*ResponseWriter).Write(startContent)
		}
		Mutex.Unlock()
	}
}

func StartTimer() {
	newTicker := time.NewTicker(time.Millisecond * 300)

	for {
		<-newTicker.C
		OnTimerTick()
	}
}
