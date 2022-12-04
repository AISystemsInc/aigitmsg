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
	out, err := exec.Command("git", "diff", "--cached").Output()
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
		Model:       "text-davinci-003",
		MaxTokens:   500,
		Prompt:      prompt,
		Temperature: 0.7,
	}

	c := gogpt.NewClient(*gptAPIKey)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)

	defer cancel()

	resp, err := c.CreateCompletion(ctx, req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(strings.TrimSpace(resp.Choices[0].Text))
}

func buildPrompt(gitDiff, gitBranch, gitTemplate string) string {
	var prompt = `write a git commit message using the following template:

`

	if gitTemplate != "" {
		prompt += gitTemplate
	} else {
		prompt += `[Action]: [Summary of the most significant change]

[bullet list of diff Highlights]

[if branch name has ticket number, display it, otherwise leave blank]

[Action] is an imperative word like Add, Remove, Delete, Fix, Refactor, Update, etc.
[Summary] is a short description of the most significant change.
[diff Highlights] is a list of the most significant changes in the diff.
`
	}

	if gitBranch != "" {
		prompt += `
The name of the branch is ` + gitBranch
	}

	if gitDiff != "" {
		prompt += `
This is the git diff:

` + gitDiff
	}

	prompt += `
Primarily consider lines that start with a + sign.`

	return prompt
}
