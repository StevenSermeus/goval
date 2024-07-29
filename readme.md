# GoVal

A simple key value database implemented in Go.

Keys are case unsensitive.

## Features

- [x] Set a key value pair
- [x] Get a value by key
- [x] Delete a key value pair
- [x] Update a key value pair
- [ ] Health check
- [x] Persist data to disk
- [x] Cache data in memory
  - [ ] Support cache invalidation and cache eviction
  - [ ] Better cache suppression
- [ ] Support ttl (Time To Live) for keys
- [x] Support remote access
  - [ ] Create a client library
- [ ] Support authentication
- [ ] Support encryption

## Network Protocol

The server uses TCP to communicate with clients. The server listens on port 8080.

All messages start with a character that represents the type of the message / response.

- `+` String: The message is a string
- `-` Error: The message is an error
- `!` Command: The message is a command

The protocol is a delimiter based protocol. The delimiter is the newline character `\n\r\n\r`.

## Known Issues and Limitations

- No support for concurrent access on the file system

## How to run

You can mount the server.yml file to the container to change the configuration. The path must be `/production/server.yml`.

```bash
docker build -t goval .
docker run -p 8080:8080 goval
```
