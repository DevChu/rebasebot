package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/chrisledet/rebasebot/git"
	"github.com/chrisledet/rebasebot/github"
)

func Status(w http.ResponseWriter, r *http.Request) {
	event := strings.ToLower(r.Method)

	log.Printf("http.request.%s.received: %s\n", event, r.RequestURI)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK\n")

	log.Printf("http.%s.response.sent: %d\n", event, http.StatusOK)
}

func Receive(w http.ResponseWriter, r *http.Request) {
	var githubEvent github.Event
	responseStatus := http.StatusOK
	event := strings.ToLower(r.Method)

	log.Printf("http.request.%s.received: %s\n", event, r.RequestURI)

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&githubEvent); err != nil {
		responseStatus = http.StatusBadRequest
		log.Printf("http.request.body.parse_failed: %s\n", err.Error())
	}

	if canRebase(githubEvent) {
		go func() {
			log.Println("bot.rebase.started")
			defer log.Println("bot.rebase.finished")

			repositoryPath := git.GetRepositoryPath(githubEvent.Repository.Name)
			pullRequest := github.FindPR(githubEvent.Repository, githubEvent.Issue.Number)

			if pullRequest.Number < 0 {
				return
			}

			branch := pullRequest.Head.Ref
			baseBranch := pullRequest.Base.Ref

			git.Fetch(repositoryPath)
			git.Checkout(repositoryPath, branch)
			git.Reset(repositoryPath, branch)
			git.Rebase(repositoryPath, baseBranch)
			git.Push(repositoryPath, branch)
		}()
	}

	w.WriteHeader(responseStatus)
	log.Printf("http.%s.response.sent: %d\n", event, responseStatus)
}

func canRebase(e github.Event) bool {
	return len(e.Repository.FullName) > 0
}