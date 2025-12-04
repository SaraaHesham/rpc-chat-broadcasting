# 💬 RPC-Chat-Broadcasting
A real-time distributed chat system implemented in Go using **RPC, goroutines, channels, and mutex-based client synchronization**.

---

## 🧠 Project Concept

Each connected client:
- Registers with the server using RPC.
- Exposes its own RPC endpoint so the server can push messages to it.
- Receives real-time messages broadcasted by the server.
- Does NOT receive its own message back (no self-echo).

The server:
- Maintains a synchronized list of active clients.
- Uses a dedicated **broadcaster goroutine** reading from a channel.
- Sends each message concurrently to all other clients.

---

## 🏗️ Architecture Overview

```
          +-------------------+        RPC Join        +----------------------+
          |      Client A     | ---------------------> |        Server        |
          | (Runs RPC server) |                        |   - client list      |
          +-------------------+                        |   - broadcaster      |
                 ^     |                                +----------------------+
                 |     |  RPC Receive()                        ^        |
                 |     +----------------------------------------+        |
                 |                                                       |
                 |                     RPC Broadcast                     |
                 +-------------------------------------------------------+
```

---

## 🧩 Concurrency Model

### Server uses:
- **sync.Mutex** → protects the shared `clients` map.
- **broadcastChan chan Message** → all messages flow through here.
- **broadcaster goroutine**:
  - waits for new messages on the channel
  - fans them out concurrently using goroutines

### No blocking
Sending a message never blocks the sender; broadcasting is fully asynchronous.

---

## 🛠️ How to Run the System

### 1. Start the server

```
cd server
go run server.go
```

Output should show:

```
Server running at 1234 ...
```

### 2. Start one or more clients

- Open multiple terminals:

```
cd client
go run client.go
```

- Enter different usernames when prompted.

```
Enter your name: Sara
Connected! Start typing...
```
### 3. Chat!

- Clients instantly receive:
  - Join notifications
  - Messages from others

- The sender does not see their own message (no self-echo).

---

## 🚀 Features
- Real-time message broadcasting
- Join notifications: `User [ID] joined`
- Concurrent message distribution using goroutines
- Mutex-protected client registry
- No message echo to the sender
- Clean separation between server & client
- Pure Go standard library (no dependencies)

---

## 📚 Learning Outcomes

Through this project you will learn:

* How RPC enables communication between distributed programs.
* How Go's concurrency model with goroutines and mutexes handles multiple clients efficiently.
* The difference between local and remote function calls.
* How to structure scalable client-server applications in Go.
* How to implement real-time message broadcasting without blocking the sender.
* How to safely manage shared state across multiple concurrent connections.
* How to design a simple but robust terminal-based chat system.
