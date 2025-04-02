# Redis-Internals

A lightweight Redis-compatible server implementation in Go, focused on understanding and learning Redis internals.

## Overview

Redis Internals is an educational project that implements core Redis functionality in Go. The project aims to provide insight into how Redis works under the hood by building a simplified version from scratch. It implements the Redis Serialization Protocol (RESP) and provides both synchronous and asynchronous server implementations.

## Features

- RESP (Redis Serialization Protocol) implementation
- Support for Simple Strings, Errors, Integers, Bulk Strings, and Arrays
- Encoding and decoding of RESP data types
- TCP server implementations
- Synchronous server (complete)
- Asynchronous server using epoll (in progress)

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

#### Asynchronous TCP Server (In Progress)

The asynchronous server will use the `epoll` system call to handle multiple connections efficiently without using multiple threads or processes.

## Future Work

* Complete the asynchronous server implementation
* Add support for more Redis commands
* Implement data structures (strings, lists, sets, hashes, sorted sets)
* Add persistence mechanisms

## I/O Multiplexing and Event Loops

Redis uses an event-driven architecture:

- **Single-Threaded Event Loop**: Similar to Node.js, Python AsyncIO
- **System Call Optimization**: Uses epoll (Linux), kqueue (BSD), or IOCP (Windows)
- **File Descriptor Monitoring**: Checks for I/O readiness without blocking
- **Non-Blocking I/O**: Performs I/O operations only when data is ready

## Acknowledgements

This project is inspired by the Redis database and is meant for educational purposes to understand the inner workings of Redis. The exploration of Redis internals, including its unique features, communication protocol, and performance architecture, is based on the insights and knowledge gained from Arpit Bhayani's course on Redis Internals.
