# Documentation for Go Voting Application

## Overview

This project implements a real-time voting system using Go, Redis, and WebSockets. Users can vote for options, and the vote counts update in real-time for all connected clients.

## File Structure

The code is organized into the following directories and files:

```lua
luaCopy code
govoting (root directory)
|-- models (directory)
|   |-- models.go
|
|-- redis_db (directory)
|   |-- redis_db.go
|
|-- websocket (directory)
|   |-- broadcast.go
|   |-- handle_client.go
|   |-- increment_vote_count.go
|   |-- upgrader.go
|
|-- index.html

```

## Code Documentation

### `models/models.go`

Defines data structures used in the application:

- **Vote**: Represents a vote cast by a user with a single field, `Option`.
- **VoteCount**: Represents the count of votes for a particular option with fields `Option` and `Count`.

### `redis_db/redis_db.go`

Handles Redis database connection and operations:

- **Clients**: Global Redis client variable.
- **init()**: Initializes the Redis client with the address `localhost:6379`.
- **Ping()**: Pings the Redis server and logs a success message if the connection is established.

### `websocket/upgrader.go`

Sets up a WebSocket upgrader:

- **upgrader**: Global WebSocket upgrader variable used to upgrade HTTP connections to WebSocket connections.

### `websocket/handle_client.go`

Handles incoming WebSocket connections from clients:

- **clients**: Global map that stores active WebSocket connections.
- **handleClient(conn)**: Handles a new WebSocket connection, reads JSON messages from the client, increments the vote count in Redis, and broadcasts the updated vote count to all connected clients.

### `websocket/increment_vote_count.go`

Increments the vote count for a given option in Redis:

- **incrementVoteCount(option)**: Increments the vote count for the given option in Redis using the HIncrBy command.

### `websocket/broadcast.go`

Broadcasts the updated vote count to all connected clients:

- **broadcastVoteCount()**: Retrieves the current vote count from Redis and sends a JSON message to each connected client with the updated vote count.

### `index.html`

Client-side HTML file that interacts with the WebSocket server:

- Sets up a WebSocket connection to `ws://localhost:8040/ws`.
- Defines three buttons for voting on options A, B, and C.
- Sends a JSON message to the WebSocket server with the selected option when a button is clicked.
- Updates the vote count display in real-time as it receives updates from the WebSocket server.

## Implementation and Testing

To implement and test this code, follow these steps:

1. **Install dependencies:**
    - Install Go.
    - Install Gorilla WebSocket: `go get github.com/gorilla/websocket`.
    - Install Redis client for Go: `go get github.com/redis/go-redis/v9`.
2. **Run Redis:**
    - Start a Redis server on `localhost:6379`.
3. **Run the Go server:**
    - Compile the Go code: `go build`.
    - Run the Go server: `go run main.go`.
4. **Open the client:**
    - Open `index.html` in a web browser.
5. **Test the application:**
    - Click on the voting buttons to cast votes.
    - Observe the vote count updating in real-time on the client-side.

**Note:** This implementation assumes a local Redis server running on `localhost:6379`. You may need to adjust the Redis connection settings depending on your environment.
