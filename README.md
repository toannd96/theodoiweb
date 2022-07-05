## Setup

### Install dependencies

- [golang version 1.18](https://go.dev/doc/install)
- [mongodb on ubuntu 20.04](https://www.digitalocean.com/community/tutorials/how-to-install-mongodb-on-ubuntu-20-04)
- [setup mongodb atlas](https://www.mongodb.com/developer/how-to/use-atlas-on-heroku/)

## Project structure

Referring from these repositories

- https://github.com/bxcodec/go-clean-arch
- https://github.com/golang-standards/project-layout

## Running

### Run locally

- Copy .env.example to .env and update .env file to suit your local environment
- Run

```
go run main.go
```

## Folder structure

```
.
├── configs
│   └── configs.go
├── db
│   ├── mongo.go
│   └── redis.go
├── Dockerfile
├── go.mod
├── go.sum
├── heroku.yml
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
│       └── middleware
│           └── cors.go
├── main.go
├── models
│   ├── session.go
│   ├── user.go
│   └── website.go
├── pkg
│   ├── duplicate.go
│   └── string.go
├── Procfile
├── README.md
└── web
    ├── static
    │   ├── assets
    │   │   └── img
    │   │       └── error-404-monochrome.svg
    │   ├── css
    │   │   └── styles.css
    │   └── js
    │       ├── record.js
    │       └── scripts.js
    └── templates
        ├── 401.html
        ├── 404.html
        ├── 500.html
        ├── dashboard.html
        ├── footer.html
        ├── header.html
        ├── layout_side_nav.html
        ├── layout_top_nav.html
        ├── login.html
        ├── password.html
        ├── register.html
        ├── tables.html
        ├── tracking.html
        └── video.html
```

## Deploy app to heroku
 
```
$ heroku login
$ heroku config:add TZ="Asia/Ho_Chi_Minh"
$ heroku addons:create heroku-redis:hobby-dev
$ heroku config --app nameapp

$ cd my-project/
$ git init
$ heroku git:remote -a nameapp
$ heroku stack:set container
$ git status
$ git add .
$ git commit -am "make it better"
$ git push heroku master
$ heroku ps:scale web=1
$ heroku logs --tail
```
