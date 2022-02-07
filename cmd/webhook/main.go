package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/asahnoln/mesproc/story"
	"github.com/asahnoln/mesproc/tg"
)

func loadStory() (*story.Story, error) {
	sFile, err := os.Open(os.Getenv("STORY_PATH"))
	if err != nil {
		return nil, fmt.Errorf("error opening story file: %w", err)
	}

	iFile, err := os.Open(os.Getenv("I18N_PATH"))
	if err != nil {
		return nil, fmt.Errorf("error opening i18n file: %w", err)
	}

	str, err := story.Load(sFile)
	if err != nil {
		return nil, fmt.Errorf("error loading story: %w", err)
	}

	i18n, err := story.LoadI18n(iFile)
	if err != nil {
		return nil, fmt.Errorf("error loading i18n: %w", err)
	}

	return str.I18n(i18n), nil
}

func createLogger() (*log.Logger, error) {
	f, err := os.OpenFile(os.Getenv("LOG_PATH"), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("error opening log file: %w", err)
	}

	logger := log.New(f, "", 0)
	return logger, nil
}

func dependendcies() (*story.Story, *log.Logger, error) {
	str, err := loadStory()
	if err != nil {
		return nil, nil, err
	}

	logger, err := createLogger()
	if err != nil {
		return nil, nil, err
	}

	return str, logger, nil
}

func main() {
	str, logger, err := dependendcies()
	if err != nil {
		log.Fatalf("error creating dependencies: %v", err)
	}
	th := tg.New(os.Getenv("BOT_ADDR"), str, logger)

	log.Fatalln(http.ListenAndServeTLS(
		os.Getenv("SRV_PORT"), os.Getenv("CERT_FILE"), os.Getenv("KEY_FILE"),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mux := http.NewServeMux()
			mux.Handle(os.Getenv("SRV_BOT_PATH"), th)
			mux.ServeHTTP(w, r)
		})))
}
