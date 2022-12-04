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
	out, err := exec.Command("git", "diff", "--cached", "--no-prefix", "-U20").Output()
	if err != nil {
		return "", fmt.Errorf("git diff failed: %s", err)
	}
	return string(out), nil
}

func getGitBranch() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return "", fmt.Errorf("git rev-parse failed: %s", err)
	}
	return string(out), nil
}

func getGitTemplate() (string, error) {
	out, err := exec.Command("git", "config", "--get", "commit.template").Output()
	if err != nil {
		log.Printf("could not get git commit template: %s", err)
		return "", nil
	}
	filename := strings.TrimSpace(string(out))
	if filename == "" {
		return "", nil
	}
	templateBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read git template file: %s", err)
	}
	return string(templateBytes), nil
}

const maxCompletionTokens = 500
const maxAllowedTokens = 4097

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
		MaxTokens:   maxCompletionTokens,
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
	var prompt = `Write a git commit message using the following rules:

`

	if gitTemplate != "" {
		prompt += gitTemplate
	} else {
		prompt += `The first line should start with an imperative word (e.g. Add, Remove, Refactor, etc.) followed by a colon and then be followed by a summary of the changes. A blank line should be next. The 3rd line should contain either a bulleted list of the general idea of changes made and/or a more detail description of the changes as a whole. The 4th line should be omitted if the branch name is a standard git branch (like main or master), otherwise if the branch name looks like a task/issue/ticket ID then add a line with "Relates to: BranchName"
`
	}

	prompt += `Here are some details about the changes you are committing:
`

	if gitBranch != "" {
		prompt += `
The branch name is ` + gitBranch
	}

	if len(prompt)+len(gitDiff) > maxAllowedTokens {
		log.Printf("prompt is too long, truncating to %d characters", maxAllowedTokens-len(gitDiff))
		gitDiff = gitDiff[:maxAllowedTokens-len(prompt)]
	}

	if gitDiff != "" {
		prompt += `
The result of the git diff command is:

` + gitDiff
	}

	prompt += `
Primarily consider lines that start with a + sign.`

	return prompt
}
