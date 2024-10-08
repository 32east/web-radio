package lib

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

var CurrentMusic Music
var Writers = make(map[*http.ResponseWriter]bool)
var Mutex = sync.Mutex{}

func OnTimerTick() {
	musicIsEnded := time.Now().After(CurrentMusic.StartTime.Add(CurrentMusic.Duration))

	if CurrentMusic.Name == "" || musicIsEnded {
		files, err := os.ReadDir("music")

		if err != nil {
			log.Fatal("Error opening folder: ", err)
			return
		}

		filesLength := len(files)
		fileName := ""

		for i := 0; i < 5; i++ {
			randomFile := files[rand.Intn(filesLength)]
			fileName = "music/" + randomFile.Name()

			if fileName != CurrentMusic.Name {
				break
			}
		}

		f, err := os.Open(fileName)
		defer f.Close()

		if err != nil {
			log.Fatal("Error opening file: ", err)
			return
		}

		content, err := GetFileContent(f)

		f, _ = os.Open(fileName)
		defer f.Close()
		duration := GetDuration(f)

		CurrentMusic = Music{fileName, time.Now(), duration, content, []byte{}, 0}

		fmt.Println("Музыка изменена: ", CurrentMusic.Name)
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

		CurrentMusic.LastEndPosition = endPosition

		if endPosition > contentLength {
			endPosition = contentLength
		}

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
