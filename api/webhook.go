package handler

import (
	"context"
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

	newContent := *dataFile.Content

	update, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	newContent = newContent[:len(newContent)-1] + "," + string(update) + "]"

	*dataFile.Content = newContent

	newGist, _, err := client.Gists.Edit(ctx, gistID, gist)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, *newGist.ID)

	w.WriteHeader(http.StatusOK)
}
