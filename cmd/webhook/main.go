package main

import (
	"log"
	"net/http"
	"os"

	"github.com/asahnoln/mesproc/story"
	"github.com/asahnoln/mesproc/tg"
)

func createStory() *story.Story {
	return story.New().
		Add(story.NewStep().
			ExpectGeo(43.257169, 76.924515, 50).
			Respond("Correct location, now enter sector number").
			Fail("Wrong location")).
		Add(story.NewStep().
			Expect("sector 1").
			Respond("move on to next sector").
			Fail(`Enter "sector 1"`)).
		Add(story.NewStep().
			Expect("sector 2").
			Respond(tg.PrefixAudio + "http://asabalar.kz/kazakhstan.mp3").
			Fail(`Enter "sector 2"`)).
		Add(story.NewStep().
			Expect("lulz").
			Respond("finish here").
			Fail(`Enter "lulz"`)).
		I18n(story.I18nMap{
			"ru": {
				"sector 1":               "сектор 1",
				"sector 2":               "сектор 2",
				"move on to next sector": "идите в следующий сектор",
				`Enter "sector 1"`:       `Введите "сектор 1"`,
				`Enter "sector 2"`:       `Введите "сектор 2"`,
				`Enter "lulz"`:           `Введите "lulz"`,
				"Correct location, now enter sector number": "Правильная локация, теперь введите сектор",
				"Wrong location": "Неправильная локация",
			},
			"kk": {
				"sector 1":               "сектор 1",
				"sector 2":               "сектор 2",
				"move on to next sector": "кз идите в следующий сектор",
				`Enter "sector 2"`:       `кз Введите "сектор 2"`,
				"Correct location, now enter sector number": "kz Правильная локация, теперь введите сектор",
				"Wrong location": "кз Неправильная локация",
			},
		})
}

func main() {
	th := tg.New(os.Getenv("BOT_ADDR"), createStory())

	log.Fatalln(http.ListenAndServeTLS(
		os.Getenv("SRV_PORT"), os.Getenv("CERT_FILE"), os.Getenv("KEY_FILE"),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mux := http.NewServeMux()
			mux.Handle(os.Getenv("SRV_BOT_PATH"), th)
			mux.ServeHTTP(w, r)
		})))
}
