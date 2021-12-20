package main

import (
	"log"
	"net/http"
	"os"

	"github.com/asahnoln/mesproc"
)

func main() {
	tg := mesproc.NewTgHandler(os.Getenv("BOT_ADDR"))
	srv := mesproc.NewServer(tg)

	log.Fatalln(http.ListenAndServeTLS(
		os.Getenv("SRV_PORT"), os.Getenv("CERT_FILE"), os.Getenv("KEY_FILE"),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mux := http.NewServeMux()
			mux.Handle(os.Getenv("SRV_BOT_PATH"), srv)
			mux.ServeHTTP(w, r)
		})))
}
