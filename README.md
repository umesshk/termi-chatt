# Termi-Chatt

**Termi-Chatt** is a terminal-based real-time chat application written in **Go** using **WebSockets**.

It allows users to create chat rooms, join existing rooms, and communicate with other users directly from the terminal.

The goal of this project is to explore **concurrency in Go**, **WebSocket communication**, and building interactive **CLI-based network applications**.

---

## Features

* Create chat rooms
* Join existing rooms
* Real-time messaging using WebSockets
* Simple terminal-based interface
* Concurrent message handling using Go routines and channels

---

## Project Structure

```
termi-chatt
│
├── cmd
│   ├── client
│   │   └── main.go        # Terminal chat client
│   └── server
│       └── main.go        # WebSocket server
│
├── internal
│   └── user
│       └── user.go        # Shared message structures
│
├── go.mod
├── Makefile
└── README.md
```

---

## How It Works

### Server

The server handles:

* WebSocket connections
* Room creation
* User joining
* Message broadcasting to all users in a room

Rooms are stored in memory and users are tracked through their WebSocket connections.

### Client

The client:

* Connects to the WebSocket server
* Sends user messages
* Receives server responses
* Displays chat messages in the terminal

---

## Installation

Clone the repository:

```bash
git clone https://github.com/umesshk/termi-chatt
cd termi-chatt
```

Install dependencies:

```bash
go mod tidy
```

---

## Running the Application

### Start the Server

```bash
make server
```

or

```bash
go run cmd/server/main.go
```

#### Environment variables

- **`PORT`**: server port (default: `8080`)
- **`POSTGRES_DSN`**: Postgres connection string (default: `host=localhost port=5432 user=postgres password=mypass dbname=termichatt sslmode=disable`)
- **`REDIS_ADDR`**: Redis address (optional). If unset, Redis features are disabled.
- **`REDIS_PASSWORD`**: Redis password (optional)
- **`REDIS_DB`**: Redis DB number (optional, default `0`)

Example:

```bash
export POSTGRES_DSN="host=localhost port=5432 user=postgres password=mypass dbname=termichatt sslmode=disable"
export REDIS_ADDR="localhost:6379"
go run cmd/server/main.go
```

---

### Start the Client

```bash
make client
```

or

```bash
go run cmd/client/main.go
```

---

## Usage

When the client starts, you will see:

```
[ 1 ] Create Room
[ 2 ] Join Room
[ 3 ] Exit
```

### Create Room

* Select option **1**
* A new chat room will be created
* Share the room ID with others

### Join Room

* Select option **2**
* Enter your username
* Enter the room ID
* Start chatting

---

## Example Chat

```
User Makito Created Room 1

Makito : hello
orn : hi
Makito : welcome!
```

---

## Technologies Used

* **Go**
* **Gorilla WebSocket**
* **Go routines**
* **Channels**

---

## Learning Goals

This project demonstrates:

* WebSocket communication in Go
* Concurrency using goroutines and channels
* CLI application development
* Handling multiple clients in real-time systems



