# Go Context Example

[![Go Report Card](https://goreportcard.com/badge/github.com/DucTran999/go-context-example)](https://goreportcard.com/report/github.com/DucTran999/go-context-example)
[![Go](https://img.shields.io/badge/Go-1.24.6-blue?logo=go)](https://golang.org)

- This repo we will explore the use of context in Go.
- See my post on medium: [Go Context Example](https://medium.com/@ductran999/go-context-example-9b5a0a0a0a0a)

## Installation

To install this repo, you can use the following command:

```bash
git clone https://github.com/DucTran999/go-context-example.git
```

---

## Usage

This project have branches:

- graceful-shutdown: example of graceful shutdown using context
- context-cancel: example of context cancellation using context
- context-timeout: example of context timeout using context
- context-deadline: example of context deadline using context

Chose example branch to run:

```bash
# example: git checkout graceful-shutdown
git checkout <branch_name>
```

Then run the example:

```bash
# Start the http server
go run main.go
```

```bash
# send request for test.
 curl -X PUT http://localhost:8080/article/2 -H \
    "Content-Type: application/json" \
    -d '{"content": "hello world"}' | jq
```

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---
