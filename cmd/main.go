package main

import (
	"log"
	"net/http"
	"os"

	"github.com/asahnoln/mesproc"
)

func main() {
	tg := mesproc.NewTgHandler(os.Getenv("BOT_ADDR"),
		mesproc.NewStory().
			Add(mesproc.NewStep().
				Expect("sector 1").
				Respond("move on to next sector").
				Fail("Enter `sector 1`")).
			Add(mesproc.NewStep().
				Expect("sector 2").
				Respond("go to next sector, yes, which is named 'lulz'").
				Fail("Enter `sector 2`")).
			Add(mesproc.NewStep().
				Expect("lulz").
				Respond("finish here").
				Fail("Enter `lulz`")),
	)
	srv := mesproc.NewServer(tg)

	log.Fatalln(http.ListenAndServeTLS(
		os.Getenv("SRV_PORT"), os.Getenv("CERT_FILE"), os.Getenv("KEY_FILE"),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mux := http.NewServeMux()
			mux.Handle(os.Getenv("SRV_BOT_PATH"), srv)
			mux.ServeHTTP(w, r)
		})))
}
