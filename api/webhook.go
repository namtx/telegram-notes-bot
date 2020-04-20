package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	DATA_FILENAME = "data.json"
)

type Update struct {
	ID             int     `json:"update_id"`
	Message        Message `json:"message"`
	UpdatedMessage Message `json:"edited_message"`
}

type Message struct {
	ID       int      `json:"message_id"`
	Text     string   `json:"text"`
	Entities []Entity `json:"entities"`
}

type Entity struct {
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	Type   string `json:"type"`
}

type GistContent struct {
	Messages []Message `json:"messages"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	gistID := os.Getenv("GIST_ID")

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GH_OAUTH_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	gist, _, err := client.Gists.Get(ctx, gistID)
	if err != nil {
		log.Fatal(err)
	}

	dataFile := gist.Files[DATA_FILENAME]

	var update Update
	var gistContent GistContent

	updateContent, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal([]byte(updateContent), &update)

	json.Unmarshal([]byte(*dataFile.Content), &gistContent)

	if update.UpdatedMessage.ID != 0 {
		var found bool

		for i := 0; i < len(gistContent.Messages); i++ {
			if gistContent.Messages[i].ID == update.UpdatedMessage.ID {
				found = true
				gistContent.Messages[i] = update.UpdatedMessage
			}
		}

		if !found {
			gistContent.Messages = append(gistContent.Messages, update.UpdatedMessage)
		}
	} else {
		gistContent.Messages = append(gistContent.Messages, update.Message)
	}

	newContent, err := json.Marshal(gistContent)
	if err != nil {
		log.Fatal(err)
	}

	*dataFile.Content = string(newContent)

	newGist, _, err := client.Gists.Edit(ctx, gistID, gist)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, *newGist.ID)
}
