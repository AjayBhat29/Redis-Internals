# Redis-Internals

A lightweight Redis-compatible server implementation in Go, focused on understanding and learning Redis internals.

## Overview

The project aims to provide insight into how Redis works under the hood by building a simplified version from scratch. It implements the Redis Serialization Protocol (RESP) and provides both synchronous and asynchronous server implementations.

## Features

- RESP (Redis Serialization Protocol) implementation
- Support for Simple Strings, Errors, Integers, Bulk Strings, and Arrays
- Encoding and decoding of RESP data types
- TCP server implementations
- Synchronous server
- Asynchronous server using epoll
- Core Redis command implementations:

  - `GET`, `SET`, `TTL`
  - `DEL`, `EXPIRE`
  - Auto-deletion of expired keys
- Eviction strategy support (e.g., Simple-First)

### Redis Serialization Protocol (RESP)

RESP is the communication protocol used by Redis:

- **Request-Response Protocol**: Clients send requests in RESP format, server responds in RESP
- **Data Type Encoding**:
  - Simple Strings: Start with '+', followed by string and CRLF
  - Integers: Start with ':', followed by integer and CRLF
  - Bulk Strings: Start with '$', include byte count, and are binary safe
  - Arrays: Start with '*', include element count, followed by RESP-encoded elements
  - Errors: Start with '-', followed by error message and CRLF

### Server Models

#### Synchronous TCP Server

The synchronous server handles each client connection in a blocking manner. It reads commands, processes them, and returns responses sequentially.

#### Asynchronous TCP Server

The asynchronous server will use the `epoll` system call to handle multiple connections efficiently without using multiple threads or processes.

The server implementation in `async_tcp.go` demonstrates how Redis achieves its impressive performance through efficient I/O handling. By only performing read system calls when there's data to be read, and by using epoll to monitor socket readiness, we avoid unnecessary waiting and context switching.

### How It Works

1. The server creates an epoll instance using `EpollCreate1()`
2. It registers the server socket to be monitored for incoming connections
3. The event loop continually checks for I/O-ready file descriptors using `EpollWait()`
4. When a new client connects, their socket is added to the epoll instance for monitoring
5. When a client socket is ready for reading, data is processed and a response is sent

This implementation mirrors Redis's efficient approach to handling network I/O without the complexity of multi-threading, while maintaining high performance and throughput.

### Command Support

The project now supports several essential Redis commands, mimicking their behavior:

- **GET**: Retrieve the value of a key
- **SET**: Assign a value to a key
- **TTL**: Return the remaining time to live for a key
- **DEL**: Remove a key from the store
- **EXPIRE**: Set a time-to-live on a key
- **Auto-deletion**: Expired keys are automatically purged
- **Eviction strategy**: Implements a simple-first policy when memory limits are approached

These features closely follow Redis semantics and provide a strong foundation for experimenting with more advanced Redis internals.

## I/O Multiplexing and Event Loops

Redis uses an event-driven architecture:

- **Single-Threaded Event Loop**: Similar to Node.js, Python AsyncIO
- **System Call Optimization**: Uses epoll (Linux), kqueue (BSD), or IOCP (Windows)
- **File Descriptor Monitoring**: Checks for I/O readiness without blocking
- **Non-Blocking I/O**: Performs I/O operations only when data is ready

## Future Work

* Implement data structures (strings, lists, sets, hashes, sorted sets)
* Add persistence mechanisms (AOF)
* Extend eviction strategies (LRU, LFU, etc.)

## Acknowledgements

This project is inspired by the Redis database and is meant for educational purposes to understand the inner workings of Redis. The exploration of Redis internals, including its unique features, communication protocol, and performance architecture, is based on the insights and knowledge gained from Arpit Bhayani's course on Redis Internals.
