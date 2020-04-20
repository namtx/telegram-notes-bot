package handler

import (
	"context"
	"fmt"
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
	gistId := os.Getenv("GIST_ID")

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GH_OAUTH_TOKEN")})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	gist, _, err := client.Gists.Get(ctx, gistId)
	if err != nil {
		log.Fatal(err)
	}

	updateContent := *gist.Files[DATA_FILENAME].Content

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, updateContent)
}
