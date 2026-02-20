# ai-agent-cli

A terminal-based AI coding agent built in Go. It uses the Anthropic API to understand your codebase and perform actions through tool calls — think a lightweight Claude Code you built yourself.

## What it does

- Talks to an LLM via HTTP (OpenAI-compatible API)
- Executes tool calls returned by the model (bash commands, file reads/writes, etc.)
- Runs an agent loop — keeps going until the task is complete
- Works entirely from your terminal

## Tech Stack

- **Language:** Go
- **API:** Anthropic / OpenAI-compatible REST API
- **Tools:** Bash execution, file system operations

## Getting Started

### Prerequisites

- Go 1.21+
- An Anthropic API key (or compatible LLM API key)

### Run

```bash
go build -o agent ./app
./agent
```

## Project Structure

```
app/
├── main.go      # Entry point, agent loop
└── tools.go     # Tool definitions and execution
```

## How it works

The agent sends your prompt to the LLM, receives a response with tool calls, executes those tools locally, feeds the results back to the model, and repeats until the task is done. It's the same core loop that powers tools like Claude Code and Cursor.
