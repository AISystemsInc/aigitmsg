package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/sashabaranov/go-gpt3"
)

func getGitDiff() (string, error) {
	out, err := exec.Command("git", "diff").Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func getGitBranch() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func getGitTemplate() (string, error) {
	out, err := exec.Command("git", "config", "--get", "commit.template").Output()
	if err != nil {
		return "", err
	}
	filename := strings.TrimSpace(string(out))
	if filename == "" {
		return "", nil
	}
	templateBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(templateBytes), nil
}

func main() {
	gptAPIKey := flag.String("gpt-key", os.Getenv("GPT_API_KEY"), "GPT API Key")
	flag.Parse()

	if *gptAPIKey == "" {
		fmt.Println("GPT API key is required")
		return
	}

	gitDiff, err := getGitDiff()
	if err != nil {
		log.Fatalf("%s", err)
	}

	gitBranch, err := getGitBranch()
	if err != nil {
		log.Fatalf("%s", err)
	}

	gitTemplate, err := getGitTemplate()
	if err != nil {
		log.Fatal(err)
	}

	prompt := buildPrompt(gitDiff, gitBranch, gitTemplate)

	req := gogpt.CompletionRequest{
		Model:     "text-davinci-003",
		MaxTokens: 500,
		Prompt:    prompt,
	}

	c := gogpt.NewClient(*gptAPIKey)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)

	defer cancel()

	resp, err := c.CreateCompletion(ctx, req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.Choices[0].Text)
}

func buildPrompt(gitDiff, gitBranch, gitTemplate string) string {
	var prompt = `Please write a commit message using the following template:

`

	if gitTemplate != "" {
		prompt += gitTemplate
	}

	if gitBranch != "" {
		prompt += `
Given the name of the branch is ` + gitBranch
	}

	if gitDiff != "" {
		prompt += `
Given the following diff:

` + gitDiff
	}

	return prompt
}
