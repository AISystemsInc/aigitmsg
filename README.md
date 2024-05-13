# aigitmsg (AI Git Message)

`aigitmsg` is a command line tool that helps you write meaningful commit messages.

## How it works

It works by taking in the output of the `git diff` command, the branch name, and an optional commit template,
and then using the OpenAI API to generate a meaningful commit message. The commit message will follow the rules you
specified in the commit template or, if no template is configured, the default commit message format.

## Installation

Install from one of the pre-built binaries from the [releases page](https://github.com/AISystemsInc/aigitmsg/releases).

Or, if you have Go installed, you can install from source:

```bash
go install github.com/AISystemsInc/aigitmsg/aigitmsg@latest
```

## Usage

1. Run `aigitmsg` with the necessary flags from the root of your repository:
   - `-gpt-key <api-key>`: Specify the OpenAI API key. If the `OPENAI_API_KEY` environment variable is set, this flag can be omitted.
   - `-model <model>`: (Optional) Specify the model to use. Default is `gpt-3.5-turbo-1106`.
   - `-context-lines <number>`: (Optional) Set the number of context lines to include around changed lines in the diff.
   - `-only-prompt`: (Optional) When set, only the prompt will be shown and the program will exit.
   - `-git-message-template <template>`: (Optional) Specify a git commit message template.
   - `-commit-msg-path <path>`: (Optional) Specify the path to the original commit message file to check for merge commits.

2. The tool will generate a meaningful commit message based on the provided inputs and configurations.

## Installing as a Git Hook

Installing `aigitmsg` as a git **prepare-commit-msg hook** is easy. First, make sure you have the `aigitmsg` binary 
saved in your `PATH` environment variable.

You can then copy this hook [prepare-commit-msg](./prepare-commit-msg) to the `.git/hooks` directory of the repository
you want to use it in. You will need to make sure the hook is executable by running `chmod +x .git/hooks/prepare-commit-msg`.

Be sure to export the env variable `OPENAI_API_KEY` with your GPT-3 API key, or you can modify the hook to add the key.
