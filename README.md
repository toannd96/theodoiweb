## Setup

### Install dependencies

- GoLang: version 1.17
- InfluxDB: version 2

## Project structure

Referring from these repositories

- https://github.com/bxcodec/go-clean-arch
- https://github.com/golang-standards/project-layout

## Running

### Run locally

- Copy .env.example to .env and update .env file to suit your local environment
- Run

```
go run cmd/analytics/main.go
```

## Folder structure

```
.
├── cmd
│   └── analytics
│       └── main.go
├── configs
│   └── configs.go
├── db
│   └── conn.go
├── go.mod
├── go.sum
├── internal
│   ├── app
│   │   └── session
│   │       ├── delivery.go
│   │       ├── delivery_http.go
│   │       ├── repository.go
│   │       └── usecase.go
│   └── pkg
│       ├── common
│       │   └── common.go
│       ├── log
│       │   ├── error.go
│       │   ├── panic.go
│       │   └── sentry.go
│       └── middleware
│           └── cors.go
├── models
│   ├── client.go
│   ├── event.go
│   ├── session.go
│   └── website.go
├── README.md
└── web
    ├── static
    │   └── record.js
    └── template
        ├── session_by_id.html
        └── session_list.html
```
