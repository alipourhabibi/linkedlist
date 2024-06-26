# Linked List Implementation in Go

This repository contains a simple implementation of a linked list in Go. It supports basic operations such as insert, find, get, and remove. Each of these operations is accompanied by comprehensive unit tests to ensure the correctness and stability of the implementation.

## Getting Started

### Prerequisites

Ensure you have Go installed on your system. You can verify this by running `go version` in your command line. If Go is not installed, you can download and install it from the [official Go website](https://golang.org/dl/).

### Running the Code

To run the main application which demonstrates some linked list operations, use the following command:

```bash
go run main.go

```

### Running All Tests

```bash
go test -v

```
### Run property-based test

```bash
go test -v -run 'TestLinkedListProperties'

go test -v -run 'TestLinkedListPropertiesQuick'

```

# Run only unit tests
```bash
go test -v -run TestLinkedListInsert
```

#### Running Specific Function Tests

```bash
go test -run TestLinkedListInsert -v

go test -run TestLinkedListFind -v

go test -run TestLinkedListGet  -v

go test -run TestLinkedListRemove   -v

```
