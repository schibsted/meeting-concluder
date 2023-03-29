# Hackday 2023 Meeting Concluder

Summarize meetings automatically, by transcribing the audio and sending a summary to Slack.

Currently, there are utilities available that can do all of this.

The web interface that binds all these together is a work in progress.

Some refactoring is needed: moving functionality from utilities to the concluder package.

## Utilities

* `cmd/rec/rec` was used for recording `cmd/rec/output.wav` which contains is a recording of me saying `This meeting is about creating a llama farm.`.
* `cmd/wav2mp4/wav2mp4` was used for converting `cmd/rec/output.wav` to `cmd/wav2mp4/output.mp4`.
* `cmd/audio2text/audio2text` was used for converting `cmd/wav2mp4/output.wav` to `cmd/audio2text/output.txt`.
* `cmd/conclude/conclude` was used for converting `cmd/audio2text/output.txt` to `cmd/conclude/output.txt`.
* `cmd/slackpost/slackpost` was used for posting `cmd/conclude/output.txt` to `#nmp-meeting-concluder` on Slack.

## Requirements

* The `ffmpeg` command (available in Homebrew), for the `cmd/wav2mp4/wav2mp4` utility.
* The `afplay` command, for the `cmd/play/play` utility.
* The Go compiler.
* A working microphone.
* An Slack web hook URL, either set as `SLACK_WEBHOOK_URL`, or as `slack_webhook` in `~/.config/concluder.toml`.
* An OpenAI API key, either set as `OPENAI_API_KEY` or `OPENAI_KEY`, or as `openai_api_key` in `~/.config/concluder.toml`.

## General info

* Version: 0.0.0
* License: Apache2
