package handler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	client := http.Client{}
	botToken := os.Getenv("TELE_BOT_TOKEN")
	request, err := http.NewRequest("GET", fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates", botToken), nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, string(body))
}
