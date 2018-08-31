# Project Zero Zero

This is a toy project to experiment with GoLang.

# Backend

## Prerequisite:
`go get -u github.com/nsf/termbox-go`

`go get github.com/go-redis/redis`


## Run
`go run newWorld.go` - flushes the db and regenerate the map

`go run createPlayer.go` - get the Id of the player

`go run client.go` - set the Id of player then run this

`go run coin.go` - spawns coin

## Reference
[Directory Layout](https://github.com/golang-standards/project-layout/blob/master/README.md)

# Frontend

Game client using phaser.js

## To start:
install yarn


## To run:
`yarn run server`

`yarn run assets`

## front-end TODO:
- [x] create rocks on board
- [x] collision detection from rocks
- [x] Accept state from stubbed websockets
- [x] Override state from stubbed websockets 
- [x] move camera and get new state
- [x] add in coin
- [x] send player move over stubbed websockets
- [ ] receive state from websockets
- [ ] emit event for movement over websockets
- [ ] fix bug where lots of players start showing up
