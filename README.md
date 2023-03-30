# Hackday Meeting Concluder

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

## Endpoints

The `cmd/restserver/restserver` executable is a web server that provides several endpopints.

### POST /login

Authenticate and receive a JWT token.

**Request:**

```json
{
  "username": "user",
  "password": "password"
}
```

**Response:**

```json
{
  "token": "your_jwt_token"
}
```

### GET /user

Get user information from the JWT token.

**Headers:**

```
Authorization: Bearer your_jwt_token
```

**Response:**

```json
{
  "name": "John Doe",
  "username": "user"
}
```

### POST /start

Start the audio recording.

**Query Parameters:**

* clap_detection: Enable or disable clap detection (boolean, default: true)
* duration: Duration of the recording in seconds (integer, default: 3600)

**Headers:**

```
Authorization: Bearer your_jwt_token
```

**Response:**

```json
{
  "message": "Started audio recording."
}
```

### POST /stop

Stop the audio recording.

**Headers:**

```
Authorization: Bearer your_jwt_token
```

**Response:**

```json
{
  "message": "Stopped audio recording."
}
```

### Example curl commands

**Login and get JWT token:**

```sh
curl -X POST -H "Content-Type: application/json" -d '{"username":"user","password":"password"}' http://localhost:3000/login
```

**Start recording with clap detection and 1-hour duration:**

```sh
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer your_jwt_token" "http://localhost:3000/start?clap_detection=true&duration=3600"
```

**Stop recording:**

```sh
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer your_jwt_token" http://localhost:3000/stop
```

## Requirements

* The PortAudio library installed on your system.
* Go 1.16 or higher.
* The `ffmpeg` command (available in Homebrew), for the `cmd/wav2mp4/wav2mp4` utility.
* The `afplay` command, for the `cmd/play/play` utility.
* A working microphone.
* An Slack web hook URL, either set as `SLACK_WEBHOOK_URL`, or as `slack_webhook` in `~/.config/concluder.toml`.
* An OpenAI API key, either set as `OPENAI_API_KEY` or `OPENAI_KEY`, or as `openai_api_key` in `~/.config/concluder.toml`.

## General info

* Version: 0.0.0
* License: Apache2
