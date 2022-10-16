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
$ cd ../leaderboard && source .env
$ docker compose build 
$ docker compose up -d
$ go run cmd/seed/main.go
```
5. Test API:
 - GET http://localhost:8080/leaderboard
 - POST http://localhost:8080/leaderboard/score
6. To stop and remove containers, networks, images, and volumes.
```
$ docker compose down -v
```

## Built using
Ubuntu 20.04, VS Code 1.17.2