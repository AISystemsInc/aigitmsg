package main

import "testing"

func TestPrompt_String(t *testing.T) {
	t.Run("should return a prompt", func(t *testing.T) {
		t.Run("with limited characters", func(t *testing.T) {
			var prompt = Prompt{
				MaxChars: 29,
			}
			prompt.Add(PromptSegment{
				Content:   "This is a prompt",
				Resizable: false,
			})
			prompt.Add(PromptSegment{
				Content:   "This could be a git diff",
				Resizable: true,
			})
			prompt.Add(PromptSegment{
				Content:   "end of prompt",
				Resizable: false,
			})

			expected := "This is a promptend of prompt"
			actual := prompt.String()

			if actual != expected {
				t.Errorf("expected %s, got %s", expected, actual)
			}

			if len(actual) > prompt.MaxChars {
				t.Errorf("expected %d characters, got %d", prompt.MaxChars, len(actual))
			}
		})

		t.Run("with enough characters", func(t *testing.T) {
			var prompt = Prompt{
				MaxChars: 200,
			}
			prompt.Add(PromptSegment{
				Content:   "This is a prompt",
				Resizable: false,
			})
			prompt.Add(PromptSegment{
				Content:   "This could be a git diff",
				Resizable: true,
			})
			prompt.Add(PromptSegment{
				Content:   "end of prompt",
				Resizable: false,
			})

			expected := "This is a promptThis could be a git diffend of prompt"
			actual := prompt.String()

			if actual != expected {
				t.Errorf("expected %s, got %s", expected, actual)
			}

			if len(actual) > prompt.MaxChars {
				t.Errorf("expected %d characters, got %d", prompt.MaxChars, len(actual))
			}
		})
	})
}
