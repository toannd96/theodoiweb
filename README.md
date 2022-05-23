## Setup

### Install dependencies

- GoLang: version 1.17

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
├── configs
│   └── configs.go
├── db
│   ├── mongo.go
│   └── redis.go
├── go.mod
├── go.sum
├── internal
│   ├── app
│   │   └── session
│   │       ├── delivery.go
│   │       ├── delivery_http.go
│   │       ├── repository.go
│   │       └── usecase.go
│   └── pkg
│       ├── duration
│       │   └── duration.go
│       ├── geodb
│       │   ├── geodb.go
│       │   └── GeoLite2-City.mmdb
│       ├── log
│       │   ├── error.go
│       │   ├── panic.go
│       │   └── sentry.go
│       └── middleware
│           └── cors.go
├── main.go
├── models
│   ├── session.go
│   ├── user.go
│   └── website.go
├── README.md
└── web
    ├── static
    │   └── record.js
    └── template
        ├── session_by_id.html
        └── session_list.html
```
