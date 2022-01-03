package main

import (
	"log"
	"net/http"
	"os"

	"github.com/asahnoln/mesproc"
)

func story() *mesproc.Story {
	return mesproc.NewStory().
		Add(mesproc.NewStep().
			Expect("sector 1").
			Respond("move on to next sector").
			Fail(`Enter "sector 1"`)).
		Add(mesproc.NewStep().
			Expect("sector 2").
			Respond("audio:http://asabalar.kz/kazakhstan.mp3").
			Fail(`Enter "sector 2"`)).
		Add(mesproc.NewStep().
			Expect("lulz").
			Respond("finish here").
			Fail(`Enter "lulz"`)).
		I18n(mesproc.I18nMap{
			"ru": {
				"sector 1":               "сектор 1",
				"sector 2":               "сектор 2",
				"move on to next sector": "идите в следующий сектор",
				`Enter "sector 1"`:       `Введите "сектор 1"`,
				`Enter "sector 2"`:       `Введите "сектор 2"`,
				`Enter "lulz"`:           `Введите "lulz"`,
			},
			"kk": {
				"sector 1":               "кз сектор 1",
				"sector 2":               "кз сектор 2",
				"move on to next sector": "кз идите в следующий сектор",
				`Enter "sector 1"`:       `кз Введите "сектор 1"`,
				`Enter "sector 2"`:       `кз Введите "сектор 2"`,
			},
		})
}

func main() {
	tg := mesproc.NewTgHandler(os.Getenv("BOT_ADDR"), story())
	srv := mesproc.NewServer(tg)

	log.Fatalln(http.ListenAndServeTLS(
		os.Getenv("SRV_PORT"), os.Getenv("CERT_FILE"), os.Getenv("KEY_FILE"),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mux := http.NewServeMux()
			mux.Handle(os.Getenv("SRV_BOT_PATH"), srv)
			mux.ServeHTTP(w, r)
		})))
}
