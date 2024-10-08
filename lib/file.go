package lib

import (
	"github.com/tcolgate/mp3"
	"os"
	"time"
)

func GetFileContent(f *os.File) ([]byte, error) {
	fileInfo, err := f.Stat()

	if err != nil {
		return nil, err
	}

	buffer := make([]byte, fileInfo.Size())
	_, err = f.Read(buffer)

	if err != nil {
		return nil, err
	}

	return buffer, nil
}

func GetDuration(f *os.File) time.Duration {
	var frame mp3.Frame

	skipped := 0

	d := mp3.NewDecoder(f)
	totalDuration := 0.0

	for {
		if err := d.Decode(&frame, &skipped); err != nil {
			break
		}

		totalDuration += frame.Duration().Seconds()
	}

	duration := time.Second * time.Duration(totalDuration)

	return duration
}
