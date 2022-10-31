<h1>LEADERBOARD</h1>

## Table of contents
* [General info](#general-info)
* [Technologies](#technologies)
* [Setup](#setup)
* [Built using](#built-using)

## General info
Leaderboard API for games

## Technologies
Project is created using:
* docker compose version: 2.11.2
* go version: 1.19.2
	
## Setup
To run this project:
1. Install technologies listed above
2. Open code editor of choice
3. Clone github repo mrzenoscout/leaderboard
4. Open terminal & run commands
```
$ cd ./leaderboard
$ docker compose build 
$ docker compose up -d
$ go run cmd/seed/main.go
```
5. Test API:
 - GET http://localhost:8080/leaderboard
    (request query params: name, page, limit, month, year, all-time)
 - POST http://localhost:8080/leaderboard/score
    (request json body: name, score)
   Use authorization bearer token:
   ```
   eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
   ```
6. To stop and remove containers, networks, images, and volumes.
```
$ docker compose down -v
```

## Built using
Ubuntu 20.04, VS Code 1.17.2