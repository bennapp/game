# Overview / Genre (really big defition, type of game etc)
mmo
grid / cell based (infinitely scales)
interact and *build*
persistent
sandbox
pvp
pve

# Gamplay (a little smaller, specific game 'interations')
emphasis for the origin on the world, 0,0

a player has their origin L*
each player is given an origin that is further out from 0,0 (spiral-like) L*
a player can change their origin (but only further out) L*

players can collect resources
player can use resource to build

players can build things to defend, attack, and gain more resources

the game encourage collaborating but their is no built in collabortive features

the player can die and respawn

some buildings can be destroy by player or npc (maybe some don't)
some buildings must be safely disconnect or to be protected
more costly in resources, can be better protected (walls)
projectiles

# Adapter
## tech
Web sockets

# Game Client
## tech
open source rendering libary possible web based / browser based (not)

Questions:
what happens to your stuff (resources building etc) when you die or log off

# Inspirations
minecraft
starcraft
osrs


STEPS
Server v0
---
- cell w/ mutex not subworld
- direction
- coin, building wall (2x1), HP
- Object Design - Player
---
- Architecture - View on Subworlds
- Cell Pagination - in Memory
---
- Figure out architecture for - Action Function
-- logging system <0,1, PLAYER1> destory <0,2 ROCK> -> <0,2 EMPTY>, PLAYER stone++
-- get coin, build, destroy
---
- Load resource from DB
- Write to DB
---

Client v0

Adapter v0

Server v1
- Multi player
- Dynamicically generated map
- Respawn

server notes
logging advantages / consistent action consumption advantages

not deadlocking
1. lock the minimal resources and most granual
2. lock and unlock in deterministic order (implicit game logic)
3. later on: timeout locking
4. needs to have a failure state (client as well)

future thoughts
should the player be on a vector and not grid
