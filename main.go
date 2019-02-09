package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	b, err := tb.NewBot(tb.Settings{
		Token:  os.Getenv("face_bot_token"),
		Poller: &tb.LongPoller{Timeout: time.Second * 10},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/hello", func(m *tb.Message) {
		b.Send(m.Sender, "Hello, World")
	})

	b.Handle(tb.OnDocument, func(m *tb.Message) {
		fmt.Println(m.Document.MIME)
	})

	b.Handle(tb.OnPhoto, func(m *tb.Message) {
		var photo = m.Photo.MediaFile()
		b.Download(photo, photo.FileID+".jpg")

		// time to handle thing
		processPhoto(photo.FileID + ".jpg")

		// delete file from disk after processing

		err := os.Remove(photo.FileID + ".jpg")
		if err != nil {
			log.Println(err)
		}
	})

	b.Start()
}

func processPhoto(fname string) {
	var file, err = os.OpenFile(fname, os.O_RDONLY, 0644)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	var text = make([]byte, 1024)
	for {
		_, err = file.Read(text)

		if err == io.EOF {
			break
		}

		if err != nil && err != io.EOF {
			log.Println(err)
			break
		}
	}

	fmt.Println(" ===> Reading from file")
	fmt.Println(string(text))
}
