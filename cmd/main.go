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

	log.Fatalln(http.ListenAndServe(os.Getenv("SRV_ADDR"), srv))
}
