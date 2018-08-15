# Overview / Genre
mmo

grid / cell based (infinitely scales)

interact and *build*

persistent

sandbox

pvp

pve

exploration

generative world

# Gameplay
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

projectiles

players buildings can spawn npcs to go and collect resources for them

players buildings can spawn npcs to attack other players and npcs

## Adapter tech overview
- Web sockets

### Game Client tech overview
- open source rendering library possible web based / browser based (not)

## Questions:
what happens to your stuff (resources building etc) when you die or log off

# Inspirations
- minecraft
- starcraft
- osrs


## Server v0
~~Object Design - Player~~
~~cell w/ mutex not subworld~~
- coin, building wall (2x1), HP
- Architecture - View on Subworlds
- Cell Pagination - in Memory
- Figure out architecture for - Action Function
- logging system <0,1, PLAYER1> destory <0,2 ROCK> -> <0,2 EMPTY>, PLAYER stone++, (get coin, build, destroy)
- Load resource from DB
- Write to DB

## Client v0

## Adapter v0

## Server v1
- Multi player
- Dynamicically generated map
- Respawn

### server notes
- logging advantages / consistent action consumption advantages

### not deadlocking process
1. lock the minimal resources and most granual
2. lock and unlock in deterministic order (implicit game logic)
3. later on: timeout locking
4. needs to have a failure state (client as well)

#### Future Thoughts
should the player be on a vector and not grid?
