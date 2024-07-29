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

## Known Issues and Limitations

- No support for concurrent access on the file system

## How to run

You can mount the server.yml file to the container to change the configuration. The path must be `/production/server.yml`.

```bash
docker build -t goval .
docker run -p 8080:8080 goval
```
