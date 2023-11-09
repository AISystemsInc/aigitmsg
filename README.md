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

1. Run `aigitmsg -gpt-key <api-key>` from the root of your repository
   - If you configure the `OPENAI_API_KEY` environment variable, you can omit the `-gpt-key` flag
2. Receive a meaningful commit message

You can optionally set the `-model <model>` flag to use more capable models.

Valid models are:
 - `text-davinci-003`
 - `text-davinci-002`
 - `text-davinci-001`
 - `text-curie-001`
 - `text-babbage-001`
 - `text-ada-001`
 - `davinci`
 - `curie`
 - `babbage`
 - `ada`
 - `gpt-4`
 - `gpt-4-1106-preview`
 - `gpt-4-0613`
 - `gpt-4-32k`
 - `gpt-4-32k-0613`
 - `gpt-3.5-turbo-1106` **default**
 - `gpt-3.5-turbo`
 - `gpt-3.5-turbo-0613`
 - `gpt-3.5-turbo-16k`
 - `gpt-3.5-turbo-16k-0613`

## Getting a GPT-3 API Key

In order to use the aigitmsg tool, you will need to obtain a GPT-3 API key. You can get your API key by signing up 
for an account at [OpenAI](https://openai.com/). Once you have an account, you can find your API key in the 
Developer Settings section of your account.

Once you have your API key, you can use it with the `-gpt-key` flag when running the `aigitmsg` tool.

## Installing as a Git Hook

Installing `aigitmsg` as a git **prepare-commit-msg hook** is easy. First, make sure you have the `aigitmsg` binary 
saved in your `PATH` environment variable.

You can then copy this hook [prepare-commit-msg](./prepare-commit-msg) to the `.git/hooks` directory of the repository
you want to use it in. You will need to make sure the hook is executable by running `chmod +x .git/hooks/prepare-commit-msg`.

Be sure to export the env variable `OPENAI_API_KEY` with your GPT-3 API key, or you can modify the hook to add the key.
