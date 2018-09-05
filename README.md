# Project Zero Zero

This is a toy project to experiment with GoLang.

# Backend

## Prerequisite:
`go get -u github.com/nsf/termbox-go`

`go get github.com/go-redis/redis`


## Run
`go run newWorld.go` - flushes the db and regenerate the map

`go run createPlayer.go` - get the Id of the player

`go run coin.go` - spawns coin (optional)

to run the debugger backend / terminal client:
`go run terminalClient.go` - set the Id of player then run this

or run the websockets server to send data to the front-end
`go run websocketServer.go`

## Reference
[Directory Layout](https://github.com/golang-standards/project-layout/blob/master/README.md)

# Frontend

Game client using phaser.js

## To start:
install yarn


## To run:
`yarn install`

`yarn run assets`

## VO: TODO:
- [x] create rocks on board
- [x] collision detection from rocks
- [x] Accept state from stubbed websockets
- [x] Override state from stubbed websockets 
- [x] move camera and get new state
- [x] add in coin
- [x] send player move over stubbed websockets
- [x] receive state from websockets
- [x] emit event for movement over websockets

## v1 prep refactors and fixes:
- [ ] Gracefully handle websockets disconnect
- [ ] Use redis `Inc` to handle ids
- [ ] Fix issue where move events are not sent as vectors and player 'skips / jumps over' coins
- [ ] fix bug where lots of players start showing up

