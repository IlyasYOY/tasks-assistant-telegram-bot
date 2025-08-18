# Tasks Assistant Telegram Bot  

> [!important] Disclaimer!
> A purely vibecoded project – it’s shared for fun and is **not** intended to
> be maintained or supported.

A small Telegram bot that lets you manage a personal task list with the help of
an LLM (OpenAI‑compatible API).  You type a plain‑text message, the bot sends
the current list + your new input to the model, receives an updated list and
stores it.  

## Table of Contents  

1. [Overview](#overview)  
2. [Features](#features)  
3. [Prerequisites](#prerequisites)  
4. [Installation & Building](#installation--building)  
5. [Configuration](#configuration)  
6. [Running the Bot](#running-the-bot)  
7. [Commands & Usage](#commands--usage)  
8. [Database & Migrations](#database--migrations)  
9. [Logging & Graceful Shutdown](#logging--graceful-shutdown)  
10. [Troubleshooting](#troubleshooting)  
11. [License](#license)  

## Overview  

The **Tasks Assistant** bot is a thin wrapper around the Telegram Bot API and
an OpenAI‑compatible chat model.  

* **Stateless UI** – Users interact only via plain text messages.  
* **LLM‑driven** – The bot builds a prompt that contains the previous task list
  (if any) and the new user input, sends it to the model, and stores the
  model’s response as the new list.  
* **SQLite persistence** – The full task list is stored per‑user in a tiny
  SQLite database, making the bot portable and easy to run locally.  

## Features  

| ✅ | Feature |
|---|---------|
| ✅ | `/start` – friendly greeting |
| ✅ | `/help` – usage information |
| ✅ | Plain‑text messages → new task (LLM‑generated list) |
| ✅ | Per‑user task storage in SQLite |
| ✅ | Configurable allowed user IDs (whitelisting) |
| ✅ | Graceful shutdown on SIGINT / SIGTERM |
| ✅ | Automatic database migrations with **goose** |
| ✅ | OpenAI‑compatible client (custom base URL, API key, model) |
| ✅ | Unknown command handling (`❓ I don't understand…`) |
| ✅ | Simple, testable code – each handler implements a small interface |

## Prerequisites  

* **Go** ≥ .22 (module aware)  
* **SQLite driver** – already vendored via `modernc.org/sqlite` (no external binary needed)  
* **OpenAI‑compatible API** – endpoint, key and model name (e.g., `gpt-4o-mini`)  

## Installation & Building  

You can either install the binary directly from the repository or build it locally.

```bash
# Install the latest released version (go will download, compile and place the binary in $GOPATH/bin)
go install github.com/IlyasYOY/tasks-assistant-tg-bot@latest
```

Or clone and build:

```bash
git clone https://github.com/IlyasYOY/tasks-assistant-tg-bot.git
cd tasks-assistant-tg-bot
go build -o tasks-assistant ./cmd/bot
```

The resulting executable (`tasks-assistant` or `tasks-assistant.exe`) is ready
to run.

## Configuration  

All configuration is driven by environment variables.  
Create a `.env` file (or export variables in your shell) with the following keys:

| Variable | Description | Required? |
|----------|-------------|-----------|
| `TASKS_ASSISTANT_TG_BOT_TELEGRAM_TOKEN` | Bot token obtained from BotFather | **Yes** |
| `TASKS_ASSISTANT_TG_BOT_OPEN_API_BASE_PATH` | Base URL of the OpenAI‑compatible API (e.g., `https://api.openai.com/v1`) | No (defaults to OpenAI public endpoint) |
| `TASKS_ASSISTANT_TG_BOT_OPEN_API_API_KEY` | API key for the LLM service | **Yes** |
| `TASKS_ASSISTANT_TG_BOT_OPEN_API_MODEL` | Model name to use (e.g., `gpt-4o-mini`) | **Yes** |
| `TASKS_ASSISTANT_TG_BOT_ALLOWED_USER_IDS` | Comma‑separated list of Telegram user IDs allowed to talk to the bot. If empty, **all** users are accepted. | No |
| `TASKS_ASSISTANT_TG_BOT_SQL_DSN` | SQLite DSN. `file::memory:?cache=shared` is used by default (in‑memory DB). Provide a file path for persistent storage, e.g., `tasks.db`. | No |

Example `.env`:

```dotenv
TASKS_ASSISTANT_TG_BOT_TELEGRAM_TOKEN=123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11
TASKS_ASSISTANT_TG_BOT_OPEN_API_BASE_PATH=https://api.openai.com/v1
TASKS_ASSISTANT_TG_BOT_OPEN_API_API_KEY=sk-XXXXXXXXXXXXXXXXXXXXXXXX
TASKS_ASSISTANT_TG_BOT_OPEN_API_MODEL=gpt-4o-mini
TASKS_ASSISTANT_TG_BOT_ALLOWED_USER_IDS=123456789,987654321
TASKS_ASSISTANT_TG_BOT_SQL_DSN=./tasks.db
```

Load the file before running:

```bash
export $(grep -v '^#' .env | xargs)
```

---  

## Running the Bot  

```bash
# If you installed via `go install`
tasks-assistant

# Or run the binary you built
./tasks-assistant
```

The bot will:

1. Load configuration.  
2. Open (or create) the SQLite DB.  
3. Apply any pending migrations (`goose up`).  
4. Start long‑polling Telegram for updates.  
5. Gracefully shut down on `SIGINT`/`SIGTERM`.  

You should see log output similar to:

```
2025/08/18 12:34:56 Bot authorized as MyTasksBot
2025/08/18 12:34:56 Database migrations applied successfully
2025/08/18 12:34:56 Listening for updates…
```

---  

## Commands & Usage  

| Command | Description |
|---------|-------------|
| `/start` | Greets the user and explains the bot’s purpose. |
| `/help`  | Shows a short usage guide (the same text you’re reading now). |
| *plain text* | Anything that is **not** a slash command is treated as a new task. The bot will: <br>1️⃣ Append the message to the existing list (if any) <br>2️⃣ Send the combined prompt to the LLM <br>3️⃣ Store the model’s response as the new list <br>4️⃣ Reply with the updated list. |
| Unknown command (e.g., `/foo`) | Bot replies with a friendly “I don’t understand” message. |

**Sample interaction**

```
User: /start
Bot: 👋 Hello! I'm *Tasks Assistant* – I can help you manage your tasks using AI.
     Just send me any text and I’ll treat it as a new task. I’ll always reply with the current task list.

User: Buy milk
Bot: # Tasks

1. Buy milk

User: Call about the project
Bot: # Tasks

1. Buy milk
2. Call Alice about the project
```

(The exact formatting depends on the LLM’s response.)

## Database & Migrations  

* **Schema** – `user_tasks` table (`user_id INTEGER PRIMARY KEY, tasks TEXT`).  
* **Migrations** – Stored in `internal/store/migrations`. The bot runs `store.Migrate(db, "./internal/store/migrations")` on start, so you never need to invoke Goose manually.  
* **Persistence** – By default the bot uses an in‑memory DB (lost on restart). Set `TASKS_ASSISTANT_TG_BOT_SQL_DSN` to a file path (e.g., `tasks.db`) for durable storage.  

## Logging & Graceful Shutdown  

* All errors are logged with context (`log.Printf`).  
* The bot runs a dedicated goroutine that processes updates; any handler error is logged but does not crash the process.  
* On `SIGINT`/`SIGTERM` the main routine stops receiving updates, closes the DB and exits cleanly.  

## Troubleshooting  

| Symptom | Likely Cause | Fix |
|---------|--------------|-----|
| Bot does not start, “missing Telegram token” | `TASKS_ASSISTANT_TG_BOT_TELEGRAM_TOKEN` not set or empty | Export the variable or add it to `.env`. |
| “AI request failed” messages | Invalid API key, wrong base URL, or model name | Verify `OPEN_API_KEY`, `OPEN_API_BASE_PATH`, `OPEN_API_MODEL`. |
| No tasks are persisted after restart | Using the default in‑memory DSN | Set `TASKS_ASSISTANT_TG_BOT_SQL_DSN` to a file path. |
| Bot replies “I don’t understand …” for a command you added | Command not registered in `handler.New` | Add the new handler to the `cmdMap` when constructing the `Handler`. |
| “unauthorized user” errors | `AllowedUserIDs` does not contain your Telegram user ID | Add your ID to `TASKS_ASSISTANT_TG_BOT_ALLOWED_USER_IDS` or leave the variable empty to allow everyone. |

## License  

MIT License – see the `LICENSE` file in the repository.  

---  

*Happy task‑keeping!*  
