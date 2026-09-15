# gansible

> **Source of truth**: This project is hosted on [GitHub](https://github.com/m0rthis/gansible). The GitLab repository at [gitlab.com/m0rthis/gansible](https://gitlab.com/m0rthis/gansible) is a read-only mirror.

```bash
go get github.com/m0rthis/gansible
```

A lightweight Go helper package for building custom Ansible modules that communicate using JSON over stdin/stdout.

## Features

- Structured result format compatible with Ansible expectations
- Helper functions for clean exit (`ExitJson`, `FailJson`)
- OS family normalization and validation utilities
- Simple integration into any Go-based Ansible module

## Usage

### Result Handling

```go
response := gansible.Result{
    Changed: true,
    Msg:     "Task completed",
    Data:    someOutput,
}
gansible.ExitJson(response)

