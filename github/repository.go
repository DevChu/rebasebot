package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Repository struct {
	FullName string `json:"full_name"`
	Name     string `json:"name"`
	GitUrl   string `json:"git_url"`
	SshUrl   string `json:"ssh_url"`
	Owner    User   `json:"owner"`
}

func (r Repository) FindPR(number int) PullRequest {
	var pr PullRequest

	log.Println("github.find_pr.started")

	prPath := fmt.Sprintf("/repos/%s/pulls/%d", r.FullName, number)
	request := NewGitHubRequest(prPath, "GET")
	response, err := httpClient.Do(request)

	defer response.Body.Close()

	if err != nil {
		log.Println("github.find_pr.failed error: %s", err)
		return pr
	}

	if response.StatusCode != http.StatusOK {
		log.Println("github.find_pr.failed status: ", response.StatusCode)
		return pr
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("github.find_pr.failed error:", err)
		return pr
	}

	if err := json.Unmarshal(body, &pr); err != nil {
		log.Println("github.find_pr.failed error:", err)
		return pr
	}

	log.Printf("github.find_pr.completed number: %d\n", pr.Number)

	return pr
}
