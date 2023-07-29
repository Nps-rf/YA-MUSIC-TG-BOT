# Go Telegram Music Bot

This script provides a Telegram bot that shares information about a currently playing track (pulled from a HTTP POST request) to the chat whenever a message is received. 

## Dependencies

- The Go [telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api) for interaction with the Telegram Bot API.
- The Go [godotenv](https://github.com/joho/godotenv) to load environment variables from a `.env` file.
- Standard Go libraries: `encoding/json`, `fmt`, `log`, `net/http`, `os`, `strings`, and `sync`.

## Structures

Two main structures are used:

- `ArtistInfo`: which includes `Title` and `Link` for a music artist.
- `TrackInfo`: which includes `Artists` (a slice of `ArtistInfo`), `Image`, `Title`, and `Link` for a track.

A global `lastTrackInfo` of type `TrackInfo` and a mutex `mutex` are declared to store and control the access to the latest track information.

## Functions

- `setLastTrackHandler()`: This function is a HTTP handler that receives a POST request with JSON body containing track information. After validating the request, it decodes the JSON into a `TrackInfo` instance and saves it as the `lastTrackInfo`. It uses `mutex` to ensure safe access to the `lastTrackInfo` in concurrent contexts.
- `main()`: The entry point of the application. It loads environment variables, initializes the Telegram bot, and enters the main event loop that waits for new messages in the bot's chat. On each received message, the bot responds with the current track information and a corresponding image. It then starts a HTTP server with `setLastTrackHandler()` as a handler for the `/set-last-track` endpoint.

## How it works

1. The application loads the environment variables from a `.env` file which include the bot token and the owner's name.
2. It then initializes the Telegram bot using the provided token.
3. Upon receiving a new message in the chat, the bot responds with a message containing information about the last track that was set via the `/set-last-track` HTTP endpoint. It sends two messages: one with the track details (title and artists) and another with the track's cover image.
4. The track information is updated by sending a HTTP POST request with JSON data to the `/set-last-track` endpoint.
