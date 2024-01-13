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
├── fly.toml
├── go.mod
├── go.sum
├── internal
│   ├── app
│   │   ├── auth
│   │   │   ├── repository.go
│   │   │   └── usecase.go
│   │   ├── session
│   │   │   ├── delivery.go
│   │   │   ├── delivery_http.go
│   │   │   ├── model.go
│   │   │   ├── repository.go
│   │   │   └── usecase.go
│   │   ├── user
│   │   │   ├── delivery.go
│   │   │   ├── delivery_http.go
│   │   │   ├── model.go
│   │   │   ├── repository.go
│   │   │   └── usecase.go
│   │   └── website
│   │       ├── delivery.go
│   │       ├── delivery_http.go
│   │       ├── model.go
│   │       ├── repository.go
│   │       └── usecase.go
│   └── pkg
│       ├── duration
│       │   ├── duration.go
│       │   └── duration_test.go
│       ├── geodb
│       │   ├── geodb.go
│       │   └── GeoLite2-City.mmdb
│       ├── middleware
│       │   ├── cors.go
│       │   └── jwt.go
│       ├── security
│       │   ├── access_token.go
│       │   ├── password.go
│       │   ├── password_test.go
│       │   ├── refresh_token.go
│       │   └── token.go
│       └── string
│           ├── string.go
│           └── string_test.go
├── main.go
├── README.md
└── web
    ├── static
    │   ├── assets
    │   │   └── img
    │   │       ├── error-404-monochrome.svg
    │   │       └── home_replay.png
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
        ├── heatmaps.html
        ├── home.html
        ├── layout_side_nav.html
        ├── layout_top_nav.html
        ├── login.html
        ├── not_record.html
        ├── not_record_today.html
        ├── not_website.html
        ├── profile.html
        ├── register.html
        ├── tables.html
        ├── tracking.html
        ├── video.html
        ├── website.html
        └── websites.html
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
### Key features
- Home page

![](https://github.com/dactoankmapydev/analytics-app/blob/master/picture/homepage.png)

- Sign up

![](https://github.com/dactoankmapydev/analytics-app/blob/master/picture/signup.png)

- Sign in

![](https://github.com/dactoankmapydev/analytics-app/blob/master/picture/signin.png)

- Profile

![](https://github.com/dactoankmapydev/analytics-app/blob/master/picture/profile.png)

- Add website

![](https://github.com/dactoankmapydev/analytics-app/blob/master/picture/add_website.png)

- List website

![](https://github.com/dactoankmapydev/analytics-app/blob/master/picture/list_website.png)

- Add tracking code 

![](https://github.com/dactoankmapydev/analytics-app/blob/master/picture/add_tracking_code.png)

- Generate tracking code 

![](https://github.com/dactoankmapydev/analytics-app/blob/master/picture/tracking_code.png)

- Records

![](https://github.com/dactoankmapydev/analytics-app/blob/master/picture/records.png)

- Replay

![](https://github.com/dactoankmapydev/analytics-app/blob/master/picture/replay.png)
