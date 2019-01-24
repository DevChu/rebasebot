package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type PullRequest struct {
	Body   string `json:"body"`
	State  string `json:"state"`
	Title  string `json:"title"`
	Number int    `json:"number"`
	Head   GitRef `json:"head"`
	Base   GitRef `json:"base"`
}

// PostComment posts a new comment on pull request via GitHub API
func (pr PullRequest) PostComment(msg string) (Comment, error) {
	log.Println("github.pr.comments.create.started")

	var err error
	var comment Comment

	createCommentPath := fmt.Sprintf("/repos/%s/issues/%d/comments", pr.Base.Repository.FullName, pr.Number)
	requestBodyAsBytes := []byte(fmt.Sprintf(`{"body":"%s"}`, msg))
	requestBody := ioutil.NopCloser(bytes.NewReader(requestBodyAsBytes))

	request := NewGitHubRequest(createCommentPath)
	request.Method = "POST"
	request.Header.Set("ContentLength", string(len(requestBodyAsBytes)))
	request.Body = requestBody
	response, err := httpClient.Do(request)

	var responseBodyAsBytes []byte

	if err != nil {
		log.Println("github.pr.comments.create.failed error:", err.Error())
		return comment, err
	}

	defer response.Body.Close()

	responseBodyAsBytes, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("github.pr.comments.create.failed error:", err)
		return comment, err
	}

	if response.StatusCode != http.StatusCreated {
		apiError := new(Error)
		json.Unmarshal(responseBodyAsBytes, apiError)
		log.Printf("github.pr.comments.create.failed status: %d, msg: %s \n", response.StatusCode, apiError.Message)
		return comment, err
	}

	json.Unmarshal(responseBodyAsBytes, &comment)

	log.Println("github.pr.comments.create.completed number:", pr.Number)
	return comment, nil
}

// Merge merges a pull request (Merge Button) via GitHub API
func (pr PullRequest) Merge() error {
	log.Println("github.pr.merge.started")

	var err error

	mergePullRequestPath := fmt.Sprintf("/repos/%s/pulls/%d/merge", pr.Base.Repository.FullName, pr.Number)
	request := NewGitHubRequest(mergePullRequestPath)
	request.Method = "PUT"
	response, err := httpClient.Do(request)

	if err != nil {
		log.Println("github.pr.merge.failed error:", err.Error())
		return err
	}

	defer response.Body.Close()

	responseBodyAsBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("github.pr.merge.failed error:", err)
		return err
	}

	if response.StatusCode != http.StatusOK {
		apiError := new(Error)
		json.Unmarshal(responseBodyAsBytes, apiError)
		log.Printf("github.pr.merge.failed status: %d, msg: %s \n", response.StatusCode, apiError.Message)
		return err
	}

	log.Println("github.pr.merge.completed number:", pr.Number)
	return nil
}
