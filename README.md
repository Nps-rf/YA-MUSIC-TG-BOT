# Telegram Music Bot and Yandex Music Track Collector

This project consists of a Telegram bot written in Go that shares information about the currently playing track on Yandex Music, along with a TamperMonkey userscript for Yandex Music that sends the currently playing track information to the server.

## Go Telegram Music Bot

The Go script provides a Telegram bot that shares information about a currently playing track (pulled from a HTTP POST request) to the chat whenever a message is received.

### Dependencies

- The Go [telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api) for interaction with the Telegram Bot API.
- The Go [godotenv](https://github.com/joho/godotenv) to load environment variables from a `.env` file.
- Standard Go libraries: `encoding/json`, `fmt`, `log`, `net/http`, `os`, `strings`, and `sync`.

### How it works

1. The application loads the environment variables from a `.env` file which include the bot token and the owner's name.
2. It then initializes the Telegram bot using the provided token.
3. Upon receiving a new message in the chat, the bot responds with a message containing information about the last track that was set via the `/set-last-track` HTTP endpoint. It sends two messages: one with the track details (title and artists) and another with the track's cover image.
4. The track information is updated by sending a HTTP POST request with JSON data to the `/set-last-track` endpoint.

## TamperMonkey Script for Yandex Music

This is a TamperMonkey userscript designed to collect the currently playing track information from Yandex Music and send it to the server that our Go script is running on.

### Dependencies

- [TamperMonkey](https://www.tampermonkey.net/) extension for your browser to run the userscript.

### How it works

1. The userscript runs on Yandex Music web pages (as defined by the `@match` metadata).
2. It checks for the currently playing track every 7 seconds, as set by `setInterval(checkTrack, 7000)`.
3. If a track is currently playing, it sends the track information to the server using the `GM_xmlhttpRequest` function. The server URL is "http://localhost:8080/set-last-track", and the track information is sent as JSON data in the request body.
4. The track information sent includes the title, artists, cover image URL, and track URL, which the server can then use to update the `lastTrackInfo` variable.

With both parts of the project running, users in your Telegram chat will get updates of your currently playing track on Yandex Music whenever they send a message.
