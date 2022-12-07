package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/sashabaranov/go-gpt3"
)

func getGitDiff() (string, error) {
	out, err := exec.Command("git", "diff", "--cached", "--no-prefix", "-U0").Output()
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
		return "", nil
	}

	filename := strings.TrimSpace(string(out))
	if filename == "" {
		return "", nil
	}

	templateBytes, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read git template file: %s", err)
	}

	return string(templateBytes), nil
}

const maxCompletionTokens = 500
const maxAllowedTokens = 1024
const maxPromptLength = maxAllowedTokens * 4

func main() {
	gptAPIKey := flag.String("gpt-key", os.Getenv("GPT_API_KEY"), "GPT API Key")
	onlyShowPrompt := flag.String("only-prompt", "false", "When set, only show the prompt and exit")
	gitMessageTemplate := flag.String("git-message-template", "", "Git commit message template")

	flag.Parse()

	if *gptAPIKey == "" {
		fmt.Println("-gpt-key or GPT_API_KEY environment variable is required")
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

	gitTemplate := *gitMessageTemplate

	if gitTemplate == "" {
		gitTemplate, err = getGitTemplate()
		if err != nil {
			log.Fatalf("%s", err)
		}
	}

	prompt := buildPrompt(gitDiff, gitBranch, gitTemplate)

	if *onlyShowPrompt != "false" {
		fmt.Println(prompt)
		return
	}

	req := gogpt.CompletionRequest{
		Model:       "text-davinci-003",
		MaxTokens:   maxCompletionTokens,
		Prompt:      prompt,
		Temperature: 0.9,
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
	var prompt = Prompt{
		MaxChars: maxPromptLength,
	}

	prompt.Add(PromptSegment{
		Content: `Write a git commit message using the following rules:

`,
	})

	if gitTemplate != "" {
		prompt.Add(PromptSegment{
			Content: gitTemplate,
		})
	} else {
		prompt.Add(PromptSegment{
			Content: `
The first line should start with an imperative word (e.g. Add, Remove, Refactor, etc.) followed by a present tense summary of the changes. A blank line should be next. The 3rd line should contain either a bulleted list of the general idea of changes made and/or a more detail description of the changes as a whole. The 4th line should be omitted if the branch name is a standard git branch (like main or master), otherwise if the branch name looks like a task/issue/ticket ID then add a line with "Relates to: branch name"
`,
		})
	}

	prompt.Add(PromptSegment{
		Content: `
Here are some details about the changes you are committing:
`,
	})

	if gitBranch != "" {
		prompt.Add(PromptSegment{
			Content: `
The branch name is ` + gitBranch,
		})
	}

	if gitDiff != "" {
		prompt.Add(PromptSegment{
			Content: `
Use the following git diff as a reference:

`,
		})

		prompt.Add(PromptSegment{
			Content:   gitDiff,
			Resizable: true,
		})

		prompt.Add(PromptSegment{
			Content: `
Focus mainly on lines that start with a + sign. All other lines are only for context.`,
		})
	}

	return prompt.String()
}

// PromptSegment is a struct representing a segment of a Prompt
type PromptSegment struct {
	Content   string // string to be displayed in the segment
	Resizable bool   // indicates whether the segment is resizable or not
}

// Prompt is a struct that contains a list of PromptSegments and a max number of chars
type Prompt struct {
	PromptSegments []PromptSegment // list of PromptSegments
	MaxChars       int             // maximum number of chars for the prompt
}

// Add adds a PromptSegment to the Prompt
func (p *Prompt) Add(s PromptSegment) {
	p.PromptSegments = append(p.PromptSegments, s)
}

// String returns a string representation of the Prompt with resizable PromptSegments being resized to fit the maximum chars
func (p *Prompt) String() string {
	var (
		charsUsed int
		prompt    string
	)

	for idx, segment := range p.PromptSegments {
		remainingRequiredChars := 0
		for _, s := range p.PromptSegments[idx:] {
			if !s.Resizable {
				remainingRequiredChars += len(s.Content)
			}
		}

		// Calculate total number of chars allocated so far
		totalCharsAlloc := charsUsed + remainingRequiredChars
		// If total is greater than max, resize if segment is resizable
		if (totalCharsAlloc + len(segment.Content)) > p.MaxChars {
			if segment.Resizable {
				// Calculate remaining chars and resize
				remainingChars := p.MaxChars - totalCharsAlloc

				if remainingChars > 0 {
					segment.Content = segment.Content[:remainingChars]
				} else {
					segment.Content = ""
				}
			}
		}

		// Add segment to prompt
		prompt += segment.Content
		// Increment used chars
		charsUsed += len(segment.Content)
	}

	return prompt
}
