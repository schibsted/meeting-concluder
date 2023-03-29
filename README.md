# Hackday 2023 Meeting Concluder

Summarize meetings automatically, by transcribing the audio and sending a summary to Slack.

* It comes with a web interface.
* It is currently a work in progress.

# Requirements

* The `ffmpeg` command line utility (available in Homebrew).
* The Go compiler.
* A working microphone.
* An Slack web hook URL, either set as `SLACK_WEBHOOK_URL`, or as `slack_webhook` in `~/.config/concluder.toml`.
* An OpenAI API key, either set as `OPENAI_API_KEY` or `OPENAI_KEY`, or as `openai_api_key` in `~/.config/concluder.toml`.
