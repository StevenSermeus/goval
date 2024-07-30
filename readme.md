# GoVal

A simple key value database implemented in Go.

Keys are case unsensitive.

## Features

- [x] Set a key value pair
- [x] Get a value by key
- [x] Delete a key value pair
- [x] Update a key value pair
- [x] Health check
- [x] Persist data to disk
- [x] Cache data in memory
- [x] Support remote access
- [ ] Support ttl (Time To Live) for keys
- [ ] Support authentication
- [ ] Support more data types
  - [ ] Support integers increment and decrement
  - [ ] Support lists
- [ ] Support encryption

## Network Protocol

The server uses TCP to communicate with clients. The server listens on port 8080.

All messages start with a character that represents the type of the message / response.

- `+` String: The value to store is a string
- `*` Array: The value to store is an array
- `:` Integer: The value to store is an integer
- `-` Error: The message is an error
- `!` Command: The message is a command (GET, DEL, HEALTHCHECK)

The protocol is a delimiter based protocol. The delimiter is the newline character `\n\r\n\r`.

## Known Issues and Limitations

- No support for concurrent access on the file system

## How to run

You can mount the server.yml file to the container to change the configuration. The path must be `/production/server.yml`.

```bash
docker build -t goval .
docker run -p 8080:8080 goval
```
