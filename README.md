# Engine Zero Zero

MMORTS Game Engine written in Go and Phaser.js

# back-end

## Prerequisite:
`go get -u github.com/nsf/termbox-go`

`go get github.com/go-redis/redis`


## Run
`go run new_world.go` - flushes the db and regenerate the map

`go run create_player.go` - get the Id of the player

`go run coin.go` - spawns coin (optional)

to run the debugger backend / terminal client:
`go run terminal_client.go` - set the Id of player then run this

or run the websockets server to send data to the front-end
`go run websocket_server.go`

## Tests
Tests in back-end are intended to be run individually for the time being. For now, we are only using tests to debug and print
usages of functions as we 'test-drive' the development of them. In the future, as things become more solidified, we may add actual tests assertions.
For now they are just debugging / printing results of functions. Tests within their respective packages are intended to be more white-box / unit like whereas
tests under tests under `tests/` are black-box integration-like tests. https://stackoverflow.com/questions/19998250/proper-package-naming-for-testing-with-the-go-language

Eg.
`go test cell/cell_test.go -v`

or
`go test tests/loader_and_saver_test.go -v`


# front-end

Game client using Phaser.js

## To start:
install yarn


## To run:
`yarn install`

`yarn run assets`


----
# Notes and Plans

# Overview / Genre
mmo 'rts-like'

grid / cell based ('infinitely' scales)

interact and *build*

persistent

sandbox

pvp

pve

exploration

generative world


# Gameplay ideas to support
emphasis for the origin on the world, 0,0

a player has their origin

each player is given an origin that is further out from 0,0 (spiral-like)

a player can change their origin (but only further out)

players can collect resources

player can use resource to build

players can build things to defend, attack, and gain more resources

the game encourage collaborating but their is no built in collabortive features

the player can die and respawn

some buildings can be destroy by player or npc (maybe some don't)

some buildings must be safely disconnect or to be protected

more costly in resources, can be better protected (walls)


## building buildings that create npcs
- macro control / no micro control
- players buildings can spawn npcs to go and collect resources for them
- players buildings can spawn npcs to attack other players and npcs
- defence has the advantage in terms of protecting your base


## logging out and persisting players buildings and etc
- packing everything up into 'you box'
    - unpack, and pack timeouts
    - logging out puts everything into an undestroyable box, with you gamertag on it
    - handle lossed connection during packing out phase


## Adapter tech overview
- Web sockets


### Game Client tech overview
- open source rendering library possible web based / browser based (not)


## Questions:
what happens to your stuff (resources building etc) when you die or log off: see logging out section


# Inspirations
- minecraft
- starcraft
- warcraft 2
- osrs
- factorio


### server notes
- logging advantages / consistent action consumption advantages


### not deadlocking process
1. lock the minimal resources and most granual
2. lock and unlock in deterministic order (implicit game logic)
3. later on: timeout locking
4. needs to have a failure state (client as well)


#### Future Thoughts and Questions
- Should the player be on a vector / floating point and not an integer grid?


#### Scaling Plans
- Use redis cluster over redis server. How do we handle event consistency with redis cluster? And handle 'dropped events'?
    - How do we handle distributed locking?
- Use kubernetes to deploy backend services. Where there is a service for different aspects of the game, spawnCoins.go, spawnSnakes.go, variousNpcs.go, websocketServer.go etc.

----

# V0.0

## Server
- [X] Object Design - Player
- [X] cell w/ mutex not subworld
- [X] coin, building wall (2x1)
- [X] HP
- [X] player limited view
- [X] storing state in redis DB
- [X] remove subworlds and grids
- [X] convert data store to json

## Client
- [X] create rocks on board
- [X] collision detection from rocks
- [X] Accept state from stubbed websockets
- [X] Override state from stubbed websockets
- [X] move camera and get new state
- [X] add in coin
- [X] send player move over stubbed websockets
- [X] receive state from websockets
- [X] emit event for movement over websockets

## Adapter
- [X] websockets

----

# V0.1
Architecutre diagram:
https://docs.google.com/drawings/d/1KoQpRLkz38vh3UNjf-xKUZWU0z8rbUOcvMJSpSRLSiQ/edit?usp=sharing

## Prep Refactors and Fixes:
- [X] Gracefully handle websockets disconnect
- [X] Use uuid for ids
- [X] Separate objects and paint
- [X] introduce struct like 'cell' which contains objects and paint
- [X] Add loader and saver services
- [X] Add items layer for items on the ground

- [ ] Add event system: player object id based event channels in redis pub/sub when event happens, propagates to all player channels around the event

- [ ] Implement player as object remove coin as object
- [ ] Can we remove setting type in the loaders? investigate this, run tests etc.
- [ ] Re-work movement, introduce 'speed limit'.
- [ ] Fix issue where move events are not sent as vectors and player 'skips / jumps over' coins
- [ ] fix bug where lots of players start showing up
- [ ] break apart basic.go into various files / packages for spawn coin and player interactions / actor movement


## Game Overview
- Dynamicically generated map
- Attacking / defending
- Building
- Respawn, and setting origin
- Multi player
- Logging out

## Server

## Client

## Adapter
